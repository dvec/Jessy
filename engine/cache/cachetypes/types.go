package cachetypes

import (
	"io/ioutil"
	"log"
	"encoding/xml"
	"sync"
)

type Help struct {
	XMLName xml.Name `xml:"category"`
	Name string `xml:"name,attr"`
	Title string `xml:"title"`
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
	Path string
	Data XMLHelp
}

func (cache *HelpCache) InitCache() {
	text, err := ioutil.ReadFile(cache.Path)
	if err != nil {
		log.Println("[ERROR] [main::engine::cache::cachetypes::types.go::HelpCache.InitCache] Failed to read file: ", err)
	}

	cache.Lock()
	defer cache.Unlock()
	if err := xml.Unmarshal(text, &cache.Data); err != nil {
		log.Println("[ERROR] [main::engine::cache::cachetypes::types.go::HelpCache.InitCache] Failed to unmarshal data: ", err)
	}
}