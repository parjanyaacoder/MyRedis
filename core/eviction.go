package core

import (
	"MyRedis/config"
)

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
	}
}
