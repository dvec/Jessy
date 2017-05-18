package cachetypes

import (
	"encoding/xml"
	"sync"
	"io/ioutil"
	"log"
)

type City struct {
	XMLName xml.Name `xml:"city"`
	Name string `xml:"name,attr"`
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

func (citiesCache *CitiesCache) InitCache(path string) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::CitiesCache.InitCache] Failed to read file: ", err)
	}

	citiesCache.Lock()
	defer citiesCache.Unlock()
	if err := xml.Unmarshal(text, &citiesCache.Data); err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::CitiesCache.InitCache] Failed to unmarshal data: ", err)
	}
}
