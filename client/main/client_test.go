package main

import (
	"fmt"
	"kvstore-server/base"
	"kvstore-server/server/api"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

var generator = base.NewRandGenerator(rand.New(rand.NewSource(time.Now().Unix())))

func fillData(t *testing.T, m map[int64]int32, size int) {
	for i := 0; i < size; i++ {
		key := generator.RandomKey()
		val := generator.RandomVal()

		_, ok1 := m[key]
		if !ok1 {
			m[key] = val
		}

		ok2 := api.Insert(key, val)
		if ok1 != !ok2 {
			t.Fail()
			fmt.Println("failed")
			return
		}
	}
}

func TestInsert(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var m = make(map[int64]int32)
	fillData(t, m, 100000)
}

func TestGet(t *testing.T) {

	var m = make(map[int64]int32)
	rand.Seed(time.Now().Unix())
	fillData(t, m, 100000)

	for i := 0; i < 100000; i++ {
		key := generator.RandomKey()

		v1, ok1 := m[key]
		v2 := api.Get(key)
		if !ok1 && v2 != -1 || ok1 && v1 != v2 {
			t.Fail()
			fmt.Println("failed")
			return
		}
	}
}

func TestUpdate(t *testing.T) {

	var m = make(map[int64]int32)
	rand.Seed(time.Now().Unix())
	fillData(t, m, 100000)
	for i := 0; i < 100000; i++ {
		key := generator.RandomKey()
		val := generator.RandomVal()

		_, ok1 := m[key]
		m[key] = val

		if ok1 != api.Update(key, val) {
			t.Fail()
			fmt.Println("failed")
			return
		}
		if !ok1 && api.Get(key) != -1 || ok1 && api.Get(key) != m[key] {
			t.Fail()
			fmt.Println("failed")
			return
		}

	}

}

func TestDelete(t *testing.T) {
	var m = make(map[int64]int32)
	rand.Seed(time.Now().Unix())
	fillData(t, m, 100000)

	for i := 0; i < 100000; i++ {
		key := generator.RandomKey()

		_, ok1 := m[key]
		delete(m, key)
		ok2 := api.Delete(key)
		if ok1 != ok2 || api.Get(key) != -1 {
			t.Fail()
			fmt.Println("failed")
			return
		}
	}
}

func TestRange(t *testing.T) {

	var m = make(map[int64]int32)
	rand.Seed(123123)

	var from, to int64 = -1000000, 1000000

	for i := 0; i < 10000; i++ {
		key := generator.RandomKeyByRange(from, to)
		val := generator.RandomValByRange(int32(from), int32(to))
		_, ok1 := m[key]
		if !ok1 {
			m[key] = val
		}

		ok2 := api.Insert(key, val)
		if ok1 != !ok2 {
			t.Fail()
			fmt.Println("failed")
			return
		}
	}

	for k := 0; k < 10; k++ {
		begin, end := generator.RandomBeginEndByRange(from, to)
		fmt.Println(begin, end)
		res1 := make([]int32, 0)
		for k, v := range m {
			if k >= begin && k <= end {
				res1 = append(res1, v)
			}
		}
		res2 := api.Range(begin, end)

		sort.Slice(res1, func(i, j int) bool {
			return res1[i] < res1[j]
		})
		sort.Slice(res2, func(i, j int) bool {
			return res2[i] < res2[j]
		})
		if len(res1) != len(res2) {
			t.Fail()
			fmt.Println("length not eq")
			fmt.Println(len(res1))
			fmt.Println(len(res2))
			return
		}
		for j := 0; j < len(res1); j++ {
			if res1[j] != res2[j] {
				t.Fail()
				fmt.Println("val not eq")
				fmt.Println(res1)
				fmt.Println(res2)
				return
			}
		}
	}
}

func fillDataParallel(t *testing.T, m *sync.Map) {
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 100; i++ {
		go func() {
			generator := base.NewRandGenerator(rand.New(rand.NewSource(time.Now().Unix())))
			for i := 0; i < 1000; i++ {
				key := generator.RandomKey()
				val := generator.RandomVal()
				_, ok1 := m.Load(key)
				m.Store(key, val)
				ok2 := api.Insert(key, val)
				if ok1 != !ok2 {
					t.Fail()
					fmt.Println("failed")
					fmt.Println(ok1, ok2, key, val)
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestInsertParallel(t *testing.T) {
	//rand.Seed(time.Now().Unix())
	var m = sync.Map{}
	fillDataParallel(t, &m)
}

func TestGetParallel(t *testing.T) {
	var m = sync.Map{}
	fillDataParallel(t, &m)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			generator := base.NewRandGenerator(rand.New(rand.NewSource(time.Now().Unix())))
			for j := 0; j < 1000; j++ {
				key := generator.RandomKey()
				v1, ok1 := m.Load(key)
				v2 := api.Get(key)
				if ok1 && v1 != v2 || !ok1 && v2 != -1 {
					t.Fail()
					return
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
