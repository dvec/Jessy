package cachetypes

import (
	"encoding/xml"
	"sync"
	"io/ioutil"
	"log"
)

type City struct {
	XMLName xml.Name `xml:"city"`
	Name 	string 	 `xml:"name,attr"` //Name of the city
}

type XMLCities struct {
	XMLName    xml.Name `xml:"xml"`
	CitiesList []City   `xml:"data>city"` //Cities array
}

//Cache for the cities game
type CitiesCache struct {
	sync.Mutex
	path string //Path to file with XML data that will be parsed
	Data XMLCities //Cache
}

//Cities cache initialization
func (citiesCache *CitiesCache) InitCache(path string) {
	//Reading XML data
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::CitiesCache.InitCache] Failed to read file: ", err)
	}

	citiesCache.Lock() //Cache locking
	defer citiesCache.Unlock()

	//Parsing XML data
	if err := xml.Unmarshal(text, &citiesCache.Data); err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::CitiesCache.InitCache] Failed to unmarshal data: ", err)
	}
}
