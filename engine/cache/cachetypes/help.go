package cachetypes

import (
	"encoding/xml"
	"sync"
	"io/ioutil"
	"log"
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
	xml.Name `xml:"xml"`
	HelpList []Help `xml:"help>category"` //TODO FIX
}

type HelpCache struct {
	sync.Mutex
	path string
	Data XMLHelp
}

func (helpCache *HelpCache) InitCache(path string) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[ERROR] [main::engine::mapCache::cachetypes::cachetypes.go::HelpCache.InitCache] Failed to read file: ", err)
	}

	helpCache.Lock()
	defer helpCache.Unlock()
	if err := xml.Unmarshal(text, &helpCache.Data); err != nil {
		log.Println("[ERROR] [main::engine::mapCache::cachetypes::cachetypes.go::HelpCache.InitCache] Failed to unmarshal data: ", err)
	}
}