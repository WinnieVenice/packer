package util

import (
	"sync"
	"time"
)

type LocalCache struct {
	cache    sync.Map
	interval int
}

type Item struct {
	val        interface{}
	createTime time.Time
	duration   time.Duration
}

func (it *Item) isExpire() bool {
	if it.duration == 0 {
		return false
	}
	return time.Since(it.createTime) > it.duration
}

func (it *Item) Value() interface{} {
	return it.val
}

func NewLocalCache(interval int) *LocalCache {
	localCache := LocalCache{
		cache:    sync.Map{},
		interval: interval,
	}
	go localCache.AutoGC()
	return &localCache
}

func (lc *LocalCache) Get(key interface{}) (interface{}, bool) {
	if v, ok := lc.cache.Load(key); ok {
		return v.(*Item).val, true
	}
	return nil, false
}

func (lc *LocalCache) Set(key, value interface{}, durationSecond int) {
	lc.cache.Store(key, &Item{
		val:        value,
		createTime: time.Now(),
		duration:   time.Duration(durationSecond * 1000 * 1000 * 1000),
	})
}

func (lc *LocalCache) Del(key interface{}) {
	lc.cache.Delete(key)
}

func (lc *LocalCache) AutoGC() {
	for {
		<-time.After(time.Duration(lc.interval) * time.Second)
		if keys := lc.getExpireKeys(); len(keys) != 0 {
			for _, key := range keys {
				lc.cache.Delete(key)
			}
		}
	}
}

func (lc *LocalCache) getExpireKeys() []interface{} {
	expireKeys := []interface{}{}
	lc.cache.Range(func(key, value interface{}) bool {
		if value.(*Item).isExpire() {
			expireKeys = append(expireKeys, key)
		}
		return true
	})
	return expireKeys
}

func (lc *LocalCache) GetAllKV() map[interface{}]interface{} {
	allKVList := map[interface{}]interface{}{}
	lc.cache.Range(func(key, value interface{}) bool {
		allKVList[key] = value
		return true
	})
	return allKVList
}
