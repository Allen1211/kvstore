package api

import (
	"kvstore-server/server/store"
	"kvstore-server/server/store/hashmap"
)

var dataStore store.Store

func init() {
	dataStore = hashmap.NewShardMap(32)
}

func Insert(key int64, val int32) bool {
	return dataStore.Insert(key, val)
}

func Update(key int64, val int32) bool {
	return dataStore.Update(key, val)
}

func Delete(key int64) bool {
	return dataStore.Delete(key)
}

func Get(key int64) int32 {
	return dataStore.Get(key)
}

func Range(begin int64, end int64) []int32 {
	return dataStore.Range(begin, end)
}
