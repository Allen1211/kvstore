package base

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

var generator = NewRandGenerator(rand.New(rand.NewSource(time.Now().Unix())))

func TestRandomKey(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100000; i++ {
		key := generator.RandomKey()
		if key > math.MaxInt64 || key < math.MinInt64 {
			t.Fail()
			return
		}
		//fmt.Println(key)
	}
}

func TestRandomKeyByRange(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var from int64 = 1
	var to int64 = 20000000
	for i := 0; i < 100000; i++ {
		key := generator.RandomKeyByRange(from, to)
		if key < from || key > to {
			t.Fail()
			return
		}
		//fmt.Println(key)
	}
}

func TestRandomVal(t *testing.T) {
	rand.Seed(2021 + int64(rand.Intn(1000)))
	var begin int32 = 1
	var end int32 = 1000000000
	for i := 0; i < 100000; i++ {
		val := generator.RandomValByRange(begin, end)
		if val > end || val < begin {
			t.Fail()
			return
		}
		//fmt.Println(val)
	}
}

func TestRandomOp(t *testing.T) {
	rand.Seed(1202 + int64(rand.Intn(1000)))
	for i := 0; i < 100000; i++ {
		val := generator.RandomOp()
		if val > NumOfOp || val < 0 {
			t.Fail()
			return
		}
		//fmt.Println(val)
	}
}

func TestRandomBeginEnd(t *testing.T) {
	rand.Seed(100 + int64(rand.Intn(1000)))
	for i := 0; i < 100000; i++ {
		begin, end := generator.RandomBeginEnd()
		if begin < math.MinInt64 || begin > math.MaxInt64 || end < math.MinInt64 || end > math.MaxInt64 ||
			begin > end {
			t.Fail()
			return
		}
		//fmt.Println(begin, end)
	}
}

func TestRandomBeginEndByRange(t *testing.T) {
	rand.Seed(100 + int64(rand.Intn(1000)))
	var from, to int64 = 1, 20000000
	for i := 0; i < 100000; i++ {
		begin, end := generator.RandomBeginEndByRange(from, to)
		if begin < from || begin > to || end < from || end > to ||
			begin > end {
			t.Fail()
			return
		}
		//fmt.Println(begin, end)
	}
}

func TestRandomOpByProbability(t *testing.T) {
	producer := NewRandomOpProducer(4, 0, 1, 2, 1)

	n := 20000000
	cnt := make([]int, NumOfOp)
	for i := 0; i < n; i++ {
		opType := producer.Next()
		cnt[opType]++
	}
	total := 0.0
	for i, c := range cnt {
		p := float64(c)/float64(n)
		fmt.Printf("%s: %f\n", OpMap[i], p)
		total += p
	}
	fmt.Printf("\ntotal: %f\n", total)
}