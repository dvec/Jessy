package cache

import (
	"strings"
	"main/conf"
)

type lineCache struct {
	Data []string
	path string
}

type mapCache struct {
	Data map[string]string
	path string
}

type dictCache struct {
	Data map[string][]string
	path string
}

type DataCache struct {
	RSSCache struct{
		News      lineCache
		Bash      lineCache
		IThappens lineCache
		Zadolbali lineCache
	}
	CommandDataCache struct{
		Help mapCache
	}
	DictionaryCache dictCache
}

func (cache *lineCache) UpdateCache() {
	 cache.Data = ParseFile(cache.path)
}

func (cache *mapCache) UpdateCache() {
	data := ParseFile(cache.path)
	newCache := map[string]string{}
	for _, entry := range data {
		entries := strings.Split(entry, "\\")
		newCache[entries[0]] = entries[1]
	}
	cache.Data = newCache
}

func (cache *dictCache) UpdateCache() {
	data := ParseFile(cache.path)
	newCache := map[string][]string{}
	for _, entry := range data {
		entries := strings.Split(entry, "\\")
		newCache[entries[0]] = entries[1:]
	}
	cache.Data = newCache
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
	cache.RSSCache.Bash.UpdateCache()
	cache.RSSCache.IThappens.UpdateCache()
	cache.RSSCache.Zadolbali.UpdateCache()
	cache.CommandDataCache.Help.UpdateCache()
	cache.DictionaryCache.UpdateCache()
}