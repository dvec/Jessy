package rss

import (
	"log"
	"encoding/xml"
	"main/kernel/cache"
	"main/kernel/cache/cachetypes"
)

type Data struct {
	XMLName  xml.Name `xml:"xml"`
	Version  string	`xml:"version,attr"`
	Encoding string	`xml:"encoding,attr"`
	Data     []string `xml:"stories>story"`
}

func UpdateRss(cache *cache.RssCaches) {
	var patches = []struct {
		cache *cachetypes.RssCache
		path string
		webPath string
	}{
		{&cache.News, "data/rss/news.xml", "http://lenta.ru/rss"},
		{&cache.Bash, "data/rss/bash.xml", "http://bash.im/rss/"},
		{&cache.IThappens, "data/rss/ithappens.xml", "http://ithappens.me/rss"},
		{&cache.Zadolbali, "data/rss/zadolbali.xml", "http://zadolba.li/rss"},
	}

	log.Print("[INFO] Start updating files")
	for _, value := range patches {
		value.cache.Lock()
		value.cache.Data = GetRSSData(value.webPath)
		value.cache.Unlock()
		log.Print("[INFO] Successfully updated ", value.path)
	}
}