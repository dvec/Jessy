package cache

import (
	"io/ioutil"
	"log"
	"strings"
)

func ParseFile(path string) []string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[ERROR] [main::engine::cache::tools.go::ParseFile()] ", err)
	}
	data := strings.Split(string(file), "\\end\\")
	for index, recording := range data {
		data[index] = strings.Trim(recording, "\n ")
	}
	if len(data) >= 2 {
		return data[:len(data)-1]
	} else {
		log.Println("[ERROR] [main::engine::cache::tools.go::ParseFile()] Bad data")
		return []string{"Error"}
	}
}