package cache

import (
	"strings"
	"main/conf"
	"sync"
)

type LineCache struct {
	sync.Mutex
	Data []string
	path string
}

type MapCache struct {
	sync.Mutex
	Data map[string]string
	path string
}

type DictCache struct {
	sync.Mutex
	Data map[string][]string
	path string
}

type DataCache struct {
	RSSCache struct{
		News      LineCache
		Bash      LineCache
		IThappens LineCache
		Zadolbali LineCache
	}
	CommandDataCache struct{
		Help MapCache
	}
	DictionaryCache DictCache
}

func (cache *LineCache) UpdateCache() {
	cache.Lock()
	cache.Data = ParseFile(cache.path)
	cache.Unlock()
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
	data := ParseFile(cache.path)
	newCache := map[string][]string{}
	for _, entry := range data {
		entries := strings.Split(entry, "\\")
		newCache[entries[0]] = entries[1:]
	}
	cache.Lock()
	cache.Data = newCache
	cache.Unlock()
}

func (cache *DataCache) InitCache() {
	cache.RSSCache.News.path = conf.RSS_DIR_PATH + "/news.dat"
	cache.RSSCache.Bash.path = conf.RSS_DIR_PATH + "/bash.dat"
	cache.RSSCache.IThappens.path = conf.RSS_DIR_PATH + "/ithappens.dat"
	cache.RSSCache.Zadolbali.path = conf.RSS_DIR_PATH + "/zadolbali.dat"
	cache.CommandDataCache.Help.path = conf.COMMANDS_DIR_PATH + "/help.dat"
	cache.DictionaryCache.path = conf.DATA_DIR_PATH + "/dictionary.bin"

	cache.UpdateCache()
}

func (cache *DataCache) UpdateCache() {
	cache.RSSCache.News.UpdateCache()
	cache.RSSCache.Bash.Lock()
	cache.RSSCache.IThappens.UpdateCache()
	cache.RSSCache.Zadolbali.UpdateCache()
	cache.CommandDataCache.Help.UpdateCache()
	cache.DictionaryCache.UpdateCache()
}