package speech

import (
	"net/url"
	"main/conf"
	"fmt"
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
	"main/web/vk"
)

//Speech API project has been frozen because vk.com doesn't allows audio uploading for groups
//TODO DEFROST
func RequestAPI(text string, kit vk.ChanKit) {
	body := url.Values{}
	body.Set("text", text)
	body.Set("format", "mp3")
	body.Set("lang", "ru")
	body.Set("speaker", conf.DefaultSpeaker)
	body.Set("emotion", conf.DefaultEmotion)
	body.Set("key", conf.YandexApiKey)

	urlPath, _ := url.ParseRequestURI(conf.YandexApiUrl)
	urlPath.Path = "/generate"
	urlStr := fmt.Sprintf("%v", urlPath)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(body.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
	}
	defer resp.Body.Close()
	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	answer := kit.MakeRequest("audio.getUploadServer", nil).Output
	fmt.Println(answer)
	if answer == nil {
		return
	}
	uploadUrl := "" //answer["upload_url"].(string)

	body = url.Values{}
	body.Set("file", string(file))

	urlPath, _ = url.ParseRequestURI(uploadUrl)
	urlStr = fmt.Sprintf("%v", urlPath)

	client = &http.Client{}
	r, _ = http.NewRequest("POST", urlStr, bytes.NewBufferString(body.Encode()))
	r.Header.Add("Content-Type", "multipart/form-data")
	resp, err = client.Do(r)
	if err != nil {
		log.Println("[ERROR] [main::web::vk::api.go] ", err)
	}
	defer resp.Body.Close()
	vkAnswer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println(string(vkAnswer))
}
