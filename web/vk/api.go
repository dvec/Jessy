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
	AccessToken 	string //Access Token to make API requests
	ChanKit		ChanKit //ChanKit for safety access to the API
}

//Executes the API request. It calls method named methodName with parameters params and returns api answer and error
//During normal running error will not be passed
func (vk *Api) Request(methodName string, params map[string]string) (map[string]interface{}, error) {
	body := url.Values{}
	body.Set("access_token", conf.VkToken)
	for paramName, param := range params {
		body.Set(paramName, param)
	}
	u, _ := url.ParseRequestURI(conf.VkApiUrl)
	u.Path = conf.VkMethodUrl + methodName
	urlStr := fmt.Sprintf("%v", u)

	//Configure parameters (URL, headers)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(body.Encode()))
	r.Header.Add("Accept", "applications/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Making the request
	resp, err := client.Do(r)
	if err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
		return nil, err
	}
	defer resp.Body.Close()

	//Reading the API answer
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//Parsing the API answer
	data := make(map[string]interface{})
	if err := json.Unmarshal(content, &data); err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
		return nil, err
	}

	return data, nil
}

//Initializes the Chan Kit
func (vk *Api) InitChanKit() {
	vk.ChanKit.AnswerChan = make(chan Answer)
	vk.ChanKit.RequestChan = make(chan Request)
}
