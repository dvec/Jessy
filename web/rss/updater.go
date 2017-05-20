package rss

import (
	"log"
	"main/kernel/cache"
	"main/kernel/cache/cachetypes"
)

//Function that updates RSS cache
func UpdateRss(cache *cache.RssCaches) {
	var patches = []struct {
		cache *cachetypes.RssCache //Data cache
		webPath string //Path to the site with RSS
	}{
		{&cache.News, "http://lenta.ru/rss"},
		{&cache.Bash, "http://bash.im/rss"},
		{&cache.IThappens, "http://ithappens.me/rss"},
		{&cache.Zadolbali, "http://zadolba.li/rss"},
	}

	log.Print("[INFO] Start updating files")
	for _, value := range patches {
		newData, err := GetRSSData(value.webPath) //Receiving RSS data from web
		if err != nil {
			log.Print("[ERROR] Failed to update RSS: ", value.webPath)
			continue
		}
		log.Print("[INFO] Successfully updated ", value.webPath)
		value.cache.Lock() //Locking cache
		value.cache.Data = newData //Updating cache
		value.cache.Unlock() //Unlocking cache
	}
}