package cachetypes

import (
	"sync"
	"main/engine/aiml"
	"log"
)

type DictCache struct {
	sync.Mutex
	Data aiml.AIML
}

func (dictCache *DictCache) UpdateCache(path string) {
	dictCache.Lock()
	if err := dictCache.Data.Learn(path); err != nil {
		log.Println("[ERROR] [main::engine::mapCache::mapCache.go::DictCache.InitCache()] Failed to update mapCache ", err)
	}
	dictCache.Unlock()
}
