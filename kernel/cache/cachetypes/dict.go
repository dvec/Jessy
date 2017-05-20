package cachetypes

import (
	"sync"
	"main/kernel/aiml"
	"log"
)

//Cache for the dictionary for the bot
type DictCache struct {
	sync.Mutex
	Data aiml.AIML //Cache data
}

func (dictCache *DictCache) UpdateCache(path string) {
	dictCache.Lock() //Cache locking
	defer dictCache.Unlock()

	//Learning AIML database
	if err := dictCache.Data.Learn(path); err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::mapCache.go::DictCache.InitCache()] Failed to update mapCache ", err)
	}
}
