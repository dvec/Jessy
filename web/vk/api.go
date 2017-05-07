package vk

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
)

const API_METHOD_URL = "https://api.vk.com/method/"

type Api struct {
	AccessToken string
	UserId		int
	ExpiresIn	int
	ChanKit		ChanKit
	debug		bool
}

func (vk *Api) Request(methodName string, params map[string]string) (map[string]interface{}, error) {
	u, err := url.Parse(API_METHOD_URL + methodName)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("access_token", vk.AccessToken)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(content, &data); err != nil {
		fmt.Println(err)
	}
	return data, nil
}

func (vk *Api) InitChanKit() {
	vk.ChanKit.AnswerChan = make(chan Answer)
	vk.ChanKit.RequestChan = make(chan Request)
}
