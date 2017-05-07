package rss

import (
	"log"
	"io/ioutil"
)

func Update()  {
	log.Print("Start updating files")
	patches := map[string][2]string{
		"news": {"data/rss/news.dat", "http://lenta.ru/rss"},
		"bash": {"data/rss/bash.dat", "http://bash.im/rss/"},
		"ithappens": {"data/rss/ithappens.dat", "http://ithappens.me/rss"},
		"zadolbali": {"data/rss/zadolbali.dat", "http://zadolba.li/rss"},
	}
	for _, value := range patches {
		bytes := []byte(ParseRss(value[1]))
		fileWriteErr := ioutil.WriteFile(value[0], bytes, 0644)
		if fileWriteErr != nil {
			log.Print("Failed to write to file: ", fileWriteErr)
		} else {
			log.Print("Successfully updated ", value[0])
		}
	}
}