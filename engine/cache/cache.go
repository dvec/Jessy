package cache

import (
	"strings"
	"main/conf"
	"sync"
	"main/engine/aiml"
	"log"
	"main/engine/cache/cachetypes"
	"math/rand"
	"main/web/rss"
)

type RssCache struct {
	sync.Mutex
	Data []string
}

type MapCache struct {
	sync.Mutex
	Data map[string]string
	path string
}

type DictCache struct {
	sync.Mutex
	Data aiml.AIML
	path string
}

type DataCache struct {
	RSSCache struct{
		News      RssCache
		Bash      RssCache
		IThappens RssCache
		Zadolbali RssCache
	}
	CommandDataCache struct{
		Help cachetypes.HelpCache
	}
	DictionaryCache DictCache
}


func (cache *RssCache) ChooseRandom() string {
	return cache.Data[rand.Intn(len(cache.Data) - 1)]
}

func (cache *MapCache) UpdateCache() {
	data := ParseFile(cache.path)
	newCache := map[string]string{}
	for _, entry := range data {
		entries := strings.Split(entry, "\\")
		newCache[entries[0]] = entries[1]
	}
	cache.Lock()
	cache.Data = newCache
	cache.Unlock()
}

func (cache *DictCache) UpdateCache() {
	cache.Lock()
	if err := cache.Data.Learn(cache.path); err != nil {
		log.Println("[ERROR] [main::engine::cache::cache.go::DictCache.InitCache()] Failed to update cache ", err)
	}
	cache.Unlock()
}

func (cache *DataCache) InitCache() {
	cache.CommandDataCache.Help.Path = conf.COMMANDS_DIR_PATH + "/help.xml"
	cache.DictionaryCache.path = conf.DATA_DIR_PATH + "/dict.aiml.xml"

	cache.DictionaryCache.Data = *aiml.NewAIML()
	cache.DictionaryCache.Data.Memory["message"] = "..."

	cache.CommandDataCache.Help.InitCache()
	cache.DictionaryCache.UpdateCache()

	cache.UpdateRssCache(rss.Update())
}

func (cache *DataCache) UpdateRssCache(newCache map[string][]string) {
	cache.RSSCache.News.Data = newCache["news"]
	cache.RSSCache.Bash.Data = newCache["bash"]
	cache.RSSCache.IThappens.Data = newCache["ithappens"]
	cache.RSSCache.Zadolbali.Data = newCache["zadolbali"]
}