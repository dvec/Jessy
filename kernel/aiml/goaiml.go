package aiml

//Code from https://github.com/eduardonunesp/goaiml

import (
	"encoding/xml"
	"main/conf"
)

type AIML struct {
	Memory map[string]string
	Bot    map[string]string
	Root   AIMLRoot
}

type AIMLRoot struct {
	XMLName    xml.Name       `xml:"aiml"`
	Categories []AIMLCategory `xml:"category"`
}

type AIMLCategory struct {
	XMLName  xml.Name     `xml:"category"`
	Pattern  AIMLPattern  `xml:"pattern"`
	Template AIMLTemplate `xml:"template"`
}

type AIMLTemplate struct {
	XMLName xml.Name `xml:"template"`
	Content string   `xml:",innerxml"`
	Looped  bool
}

type AIMLPattern struct {
	XMLName xml.Name `xml:"pattern"`
	Content string   `xml:",innerxml"`
}

func NewAIML() *AIML {
	ret := &AIML{
		Memory: make(map[string]string),
		Bot:    make(map[string]string),
	}

	ret.Bot["name"] = conf.BotName

	return ret
}