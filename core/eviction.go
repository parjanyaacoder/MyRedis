package core

import (
	"MyRedis/config"
	"time"
)

func getCurrentClock() uint32 {
	return uint32(time.Now().UnixMilli()) & 0x00FFFFFF
}

func getIdleTime(lastAccessedAt uint32) uint32 {
	c := getCurrentClock()

	if c >= lastAccessedAt {
		return c - lastAccessedAt
	}
	return (0x00FFFFFF - lastAccessedAt) + c
}

func populateEvictionPool() {
	sampleSize := 5

	for k := range store {
		ePool.Push(k, store[k].LastAccessedAt)
		sampleSize--
		if sampleSize == 0 {
			break
		}
	}
}

func evictAllKeysLRU() {
	populateEvictionPool()
	evictCount := int16(config.EvictionRatio * float64(config.KeysLimit))

	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item == nil {
			return
		}
		Del(item.key)
	}

}

func evictAllkeysRandom() {
	evictCount := int64(float64(config.KeysLimit) * float64(config.EvictionRatio))

	for key := range store {
		Del(key)
		evictCount--
		if evictCount <= 0 {
			break
		}
	}
}

func evictFirst() {
	for key := range store {
		Del(key)
		return
	}
}

func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		evictAllkeysRandom()
	case "allkeys-lru":
		evictAllKeysLRU()
	}
}
