package vk

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"log"
	"bytes"
	"main/conf"
)


type Api struct {
	AccessToken 	string
	ChanKit		ChanKit
}

func (vk *Api) Request(methodName string, params map[string]string) (map[string]interface{}, error) {
	body := url.Values{}
	body.Set("access_token", conf.TOKEN)
	for paramName, param := range params {
		body.Add(paramName, param)
	}
	u, _ := url.ParseRequestURI(conf.API_URL)
	u.Path = "/method/" + methodName
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(body.Encode()))
	r.Header.Add("Accept", "applications/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	if err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(content, &data); err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
	}
	return data, nil
}

func (vk *Api) JsonRequest(methodName string, params map[string]string) (map[string]interface{}, error) {
	data := url.Values{}
	data.Set("access_token", conf.TOKEN)
	for paramName, param := range params {
		data.Add(paramName, param)
	}
	u, _ := url.ParseRequestURI(conf.API_URL)
	u.Path = "/method/" + methodName
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Accept", "applications/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	if err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
	}
	defer resp.Body.Close()
	text, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(text))
	return nil, nil
}

func (vk *Api) InitChanKit() {
	vk.ChanKit.AnswerChan = make(chan Answer)
	vk.ChanKit.RequestChan = make(chan Request)
}
