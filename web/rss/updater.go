package rss

import (
	"log"
	"encoding/xml"
)

type Data struct {
	XMLName  xml.Name `xml:"xml"`
	Version  string	`xml:"version,attr"`
	Encoding string	`xml:"encoding,attr"`
	Data     []string `xml:"stories>story"`
}

func Update() map[string][]string {
	out := make(map[string][]string)
	log.Print("[INFO] Start updating files")
	patches := map[string][2]string{
		"news": {"data/rss/news.xml", "http://lenta.ru/rss"},
		"bash": {"data/rss/bash.xml", "http://bash.im/rss/"},
		"ithappens": {"data/rss/ithappens.xml", "http://ithappens.me/rss"},
		"zadolbali": {"data/rss/zadolbali.xml", "http://zadolba.li/rss"},
	}

	for name, value := range patches {
		stories := ParseRss(value[1])
		out[name] = stories
		log.Print("[INFO] Successfully updated ", value[0])
	}
	return out
}