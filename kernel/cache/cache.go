package cache

import (
	"main/conf"
	"main/kernel/aiml"
	"main/kernel/cache/cachetypes"
)

//Cache for RSS data
type RssCaches struct {
	News      cachetypes.RssCache
	Bash      cachetypes.RssCache
	IThappens cachetypes.RssCache
	Zadolbali cachetypes.RssCache
}

//Cache for functions data
type CommandDataCaches struct {
	Help	cachetypes.HelpCache
	Cities 	cachetypes.CitiesCache
}

//Main cache
type DataCache struct {
	RssCache RssCaches
	CommandDataCache CommandDataCaches
	DictionaryCache cachetypes.DictCache
}

//Cache initialization
func (cache *DataCache) InitCache() {
	cache.DictionaryCache.Data = *aiml.NewAIML()

	cache.CommandDataCache.Help.InitCache(conf.CommandsDirPath + "/help.xml")
	cache.CommandDataCache.Cities.InitCache(conf.CommandsDirPath + "/cities.xml")
	cache.DictionaryCache.UpdateCache(conf.DataDirPath + "/dict.aiml.xml")
}