package cache

import (
	"io/ioutil"
	"log"
	"encoding/xml"
	"sync"
	"main/engine/aiml"
	"math/rand"
)

type Help struct {
	XMLName xml.Name `xml:"category"`
	Name string `xml:"name,attr"`
	Title string `xml:"title"`
	State string `xml:"state"`
	Description string `xml:"description"`
	Samples []struct {
		XMLName xml.Name `xml:"sample"`
		Body string `xml:"body"`
		Out string `xml:"out"`
	} `xml:"samples>sample"`
}

type XMLHelp struct {
	XMLName xml.Name `xml:"xml"`
	HelpList []Help `xml:"help>category"`
}

type HelpCache struct {
	sync.Mutex
	path string
	Data XMLHelp
}

func (helpCache *HelpCache) InitCache() {
	text, err := ioutil.ReadFile(helpCache.path)
	if err != nil {
		log.Println("[ERROR] [main::engine::mapCache::types::types.go::HelpCache.InitCache] Failed to read file: ", err)
	}

	helpCache.Lock()
	defer helpCache.Unlock()
	if err := xml.Unmarshal(text, &helpCache.Data); err != nil {
		log.Println("[ERROR] [main::engine::mapCache::types::types.go::HelpCache.InitCache] Failed to unmarshal data: ", err)
	}
}

type City struct {
	XMLName xml.Name `xml:"city"`
	Name string `xml:"name"`
}

type XMLCities struct {
	XMLName xml.Name `xml:"xml"`
	CitiesList []City `xml:"data>city"`
}

type CitiesCache struct {
	sync.Mutex
	path string
	Data XMLCities
}

func (citiesCache *CitiesCache) InitCache() {
	text, err := ioutil.ReadFile(citiesCache.path)
	if err != nil {
		log.Println("[ERROR] [main::engine::mapCache::types::types.go::CitiesCache.InitCache] Failed to read file: ", err)
	}

	citiesCache.Lock()
	defer citiesCache.Unlock()
	if err := xml.Unmarshal(text, &citiesCache.Data); err != nil {
		log.Println("[ERROR] [main::engine::mapCache::types::types.go::CitiesCache.InitCache] Failed to unmarshal data: ", err)
	}
}

type RssCache struct {
	sync.Mutex
	Data []string
}

func (rssCache *RssCache) ChooseRandom() string {
	return rssCache.Data[rand.Intn(len(rssCache.Data) - 1)]
}

type DictCache struct {
	sync.Mutex
	Data aiml.AIML
	Path string
}

func (dictCache *DictCache) UpdateCache() {
	dictCache.Lock()
	if err := dictCache.Data.Learn(dictCache.Path); err != nil {
		log.Println("[ERROR] [main::engine::mapCache::mapCache.go::DictCache.InitCache()] Failed to update mapCache ", err)
	}
	dictCache.Unlock()
}