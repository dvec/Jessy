package cachetypes

import (
	"encoding/xml"
	"sync"
	"io/ioutil"
	"log"
)

type Help struct {
	XMLName	    xml.Name `xml:"category"`
	Name 	    string   `xml:"name,attr"` //Command name
	Title	    string   `xml:"title"` //Help title
	State	    string   `xml:"state"` //Command state
	Description string   `xml:"description"` //Command description
	Samples	    []struct {
		XMLName xml.Name `xml:"sample"`
		Body    string   `xml:"body"` //Sample input
		Out     string   `xml:"out"` //Sample output
	}		     `xml:"samples>sample"`
}

type XMLHelp struct {
	xml.Name 	`xml:"xml"`
	HelpList []Help `xml:"help>category"` //Help writings array
}

//Cache for the help command
type HelpCache struct {
	sync.Mutex
	path string //Path to file with XML data that will be parsed
	Data XMLHelp //Cache
}

func (helpCache *HelpCache) InitCache(path string) {
	//Reading XML
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::HelpCache.InitCache] Failed to read file: ", err)
	}

	helpCache.Lock() //Cache locking
	defer helpCache.Unlock()

	//Parsing data
	if err := xml.Unmarshal(text, &helpCache.Data); err != nil {
		log.Println("[ERROR] [main::kernel::mapCache::cachetypes::cachetypes.go::HelpCache.InitCache] Failed to unmarshal data: ", err)
	}
}