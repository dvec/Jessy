package cache

import (
	"main/conf"
	"main/engine/aiml"
	"main/engine/cache/cachetypes"
)

type RssCaches struct {
	News      cachetypes.RssCache
	Bash      cachetypes.RssCache
	IThappens cachetypes.RssCache
	Zadolbali cachetypes.RssCache
}

type CommandDataCaches struct {
	Help	cachetypes.HelpCache
	Cities 	cachetypes.CitiesCache
}

type DataCache struct {
	RssCache RssCaches
	CommandDataCache CommandDataCaches
	DictionaryCache cachetypes.DictCache
}

func (cache *DataCache) InitCache() {
	cache.DictionaryCache.Data = *aiml.NewAIML()

	cache.CommandDataCache.Help.InitCache(conf.COMMANDS_DIR_PATH + "/help.xml")
	cache.CommandDataCache.Cities.InitCache(conf.COMMANDS_DIR_PATH + "/cities.xml")
	cache.DictionaryCache.UpdateCache(conf.DATA_DIR_PATH + "/dict.aiml.xml")
}

func (cache *RssCaches) UpdateRssCache(newCache map[string][]string) {
	cache.News.Data = newCache["news"]
	cache.Bash.Data = newCache["bash"]
	cache.IThappens.Data = newCache["ithappens"]
	cache.Zadolbali.Data = newCache["zadolbali"]
}