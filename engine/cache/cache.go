package cache

import (
	"main/conf"
	"main/engine/aiml"
	"main/web/rss"
)

type DataCache struct {
	RSSCache struct{
		News      RssCache
		Bash      RssCache
		IThappens RssCache
		Zadolbali RssCache
	}
	CommandDataCache struct{
		Help	HelpCache
		Cities 	CitiesCache
	}
	DictionaryCache DictCache
}

func (cache *DataCache) InitCache() {
	cache.CommandDataCache.Help.path = conf.COMMANDS_DIR_PATH + "/help.xml"
	cache.CommandDataCache.Cities.path = conf.COMMANDS_DIR_PATH + "/cities.xml"
	cache.DictionaryCache.Path = conf.DATA_DIR_PATH + "/dict.aiml.xml"

	cache.DictionaryCache.Data = *aiml.NewAIML()
	cache.DictionaryCache.Data.Memory["message"] = "..."

	cache.CommandDataCache.Help.InitCache()
	cache.CommandDataCache.Cities.InitCache()
	cache.DictionaryCache.UpdateCache()

	cache.UpdateRssCache(rss.Update())
}

func (cache *DataCache) UpdateRssCache(newCache map[string][]string) {
	cache.RSSCache.News.Data = newCache["news"]
	cache.RSSCache.Bash.Data = newCache["bash"]
	cache.RSSCache.IThappens.Data = newCache["ithappens"]
	cache.RSSCache.Zadolbali.Data = newCache["zadolbali"]
}