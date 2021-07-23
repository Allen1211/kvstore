package hashmap

import (
	"sync"
)

const magic = 1.5

type shard struct {
	m map[int64]int32
	sync.RWMutex
}

type ShardMap struct {
	shards    []*shard
	numShards int
}

func NewShardMap(numShards int) ShardMap {
	res := new(ShardMap)
	res.numShards = numShards
	res.shards = make([]*shard, res.numShards)
	for i := 0; i < len(res.shards); i++ {
		res.shards[i] = &shard{
			m: make(map[int64]int32),
		}
	}
	return *res
}

/* implement interface store.Store */

func (sm ShardMap) Insert(key int64, val int32) bool {
	return sm.setIfAbsent(key, val)
}

func (sm ShardMap) Update(key int64, val int32) bool {
	return sm.setIfPresent(key, val)
}

func (sm ShardMap) Delete(key int64) bool {
	return sm.removeIfPresent(key)
}

func (sm ShardMap) Get(key int64) int32 {
	if res, ok := sm.get(key); ok {
		return res
	} else {
		return -1
	}
}

func (sm ShardMap) Range(begin int64, end int64) []int32 {
	var res []int32
	if end < begin {
		return res
	} else if float64(end-begin)*magic <= float64(sm.size()) {
		for k := begin; k <= end; k++ {
			if v, ok := sm.get(k); ok {
				res = append(res, v)
			}
		}
		return res
	} else {
		return sm.rangeGet(begin, end)
	}
}

/* private */

func (sm ShardMap) get(k int64) (int32, bool) {
	s := sm.getShard(k)
	s.RLock()
	res, ok := s.m[k]
	s.RUnlock()
	return res, ok
}

func (sm ShardMap) setIfAbsent(k int64, v int32) bool {
	s := sm.getShard(k)

	if s.ContainsKey(k) {
		return false
	}

	s.Lock()
	s.m[k] = v
	s.Unlock()
	return true
}

func (sm ShardMap) setIfPresent(k int64, v int32) bool {
	s := sm.getShard(k)

	if !s.ContainsKey(k) {
		return false
	}

	s.Lock()
	s.m[k] = v
	s.Unlock()
	return true
}

func (sm ShardMap) removeIfPresent(k int64) bool {
	s := sm.getShard(k)

	if !s.ContainsKey(k) {
		return false
	}

	s.Lock()
	delete(s.m, k)
	s.Unlock()
	return true
}

func (sm ShardMap) size() int {
	res := 0
	for _, s := range sm.shards {
		s.RLock()
		res += s.size()
		s.RUnlock()
	}
	return res
}

func (sm ShardMap) rangeGet(begin int64, end int64) []int32 {
	res := make([]int32, 0)
	for _, s := range sm.shards {
		s.RLock()
		for k, v := range s.m {
			if k >= begin && k <= end {
				res = append(res, v)
			}
		}
		s.RUnlock()
	}
	return res
}

func (s *shard) ContainsKey(k int64) bool {
	s.RLock()
	_, ok := s.m[k]
	s.RUnlock()
	return ok
}

func (s *shard) size() int {
	return len(s.m)
}

func (sm ShardMap) getShard(k int64) *shard {
	idx := hash64Bit(k) % sm.numShards
	return sm.shards[idx]
}

func hash64Bit(i int64) int {
	return (int)(i ^ (i >> 32))
}
