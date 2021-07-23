package hashmap

import (
	"sync"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	benchmarkGet(b, 256)
}

func BenchmarkGetRaw(b *testing.B) {
	benchmarkGetRaw(b)
}

func BenchmarkGetSync(b *testing.B) {
	benchmarkGetSync(b)
}

func benchmarkGet(b *testing.B, numOfShards int) {
	m := createShardMap(numOfShards)
	for i := 0; i < b.N; i++ {
		m.Insert(int64(i), 123)
	}
	var wg sync.WaitGroup
	wg.Add(b.N)
	getF := func(key int64) {
		for j := 0; j < 1; j++ {
			m.Get(key)
		}
		wg.Done()
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		go getF(int64(j))
	}
	wg.Wait()
}

func benchmarkGetSync(b *testing.B) {
	m := sync.Map{}
	for i := 0; i < b.N; i++ {
		m.Store(int64(i), 123)
	}
	var wg sync.WaitGroup
	wg.Add(b.N)
	getF := func(key int64) {
		for j := 0; j < 1; j++ {
			m.Load(key)
		}
		wg.Done()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go getF(int64(i))
	}
	wg.Wait()
}

func benchmarkGetRaw(b *testing.B) {
	m := make(map[int64]int32)
	for i := 0; i < b.N; i++ {
		m[int64(i)] = 123
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m[int64(i)]
	}
}

func BenchmarkSet(b *testing.B) {
	benchmarkSet(b, 256)
}

func benchmarkSet(b *testing.B, numOfShards int) {
	m := createShardMap(numOfShards)
	var wg sync.WaitGroup
	wg.Add(b.N)
	setF := func(key int64, val int32) {
		for j := 0; j < 1; j++ {
			m.Update(key, val)
		}
		wg.Done()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go setF(int64(i), 123)
	}
	wg.Wait()
}

func benchmarkFixedGetSet(b *testing.B, numOfShards int) {
	m := createShardMap(numOfShards)
	for i := 0; i < b.N; i++ {
		m.Insert(int64(i), 123)
	}

	var wg sync.WaitGroup
	wg.Add(b.N)
	getF := func(key int64) {
		for j := 0; j < 9; j++ {
			m.Get(key)
		}
		wg.Done()
	}
	setF := func(key int64, val int32) {
		for j := 0; j < 1; j++ {
			m.Update(key, val)
		}
		wg.Done()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go getF(int64(i))
		go setF(int64(i), 123)
	}
	wg.Wait()
}

func benchmarkFixedGetSetSyncMap(b *testing.B) {
	m := sync.Map{}
	//for i := 0; i < b.N; i++ {
	//	m.Store(int64(i), 123)
	//}
	//b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(2 * b.N)
	getF := func(key int64) {
		for j := 0; j < 9; j++ {
			m.Load(key)
		}
		wg.Done()
	}
	setF := func(key int64, val int32) {
		for j := 0; j < 1; j++ {
			m.Store(key, val)
		}
		wg.Done()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go getF(int64(i))
		go setF(int64(i), 123)
	}
	wg.Wait()
}

func BenchmarkFixedGetSet(b *testing.B) {
	benchmarkFixedGetSet(b, 64)
}

func BenchmarkFixedGetSetSync(b *testing.B) {
	benchmarkFixedGetSetSyncMap(b)
}

func createShardMap(numOfShards int) ShardMap {
	return NewShardMap(numOfShards)
}
