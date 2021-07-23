package store

type Key int64
type Val int32

type Store interface {
	Insert(key int64, val int32) bool
	Update(key int64, val int32) bool
	Delete(key int64) bool
	Get(key int64) int32
	Range(begin int64, end int64) []int32
}

type MultiStoreOp interface {
	MInsert(keys []int64, vals []int32) []bool
	MUpdate(keys []int64, vals []int32) []bool
	MDelete(keys []int64) []bool
	MGet(keys []int64) []int32
}
