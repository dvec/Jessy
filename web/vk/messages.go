package vk

import (
	"log"
	"net/http"
	"fmt"
	"main/conf"
	"io/ioutil"
	"encoding/json"
)

type Message struct {
	Id          int64
	UserId      int64
	Text        string
	Attachments map[string]interface{}
}

type NewMessageChan chan Message

type LongPoll struct {
	key    string
	server string
	ts     int64
	NewMessageChan
}

const (
	//FLAGS
	unread = 1  //Unread message flag
	chat   = 16 //Chat message flag

	//LONG POLL
	requestBody = "https://%v?act=a_check&key=%v&ts=%v&wait=%v&mode=2&version=1"
)

func (lp *LongPoll) Start(chanKit ChanKit) {
	answer := chanKit.MakeRequest("messages.getLongPollServer", nil)
	if answer.Error != nil {
		log.Println("[ERROR] [Messages::init]:", answer.Error)
	}
	if answer.Output["response"] == nil {
		log.Println("[ERROR]	[Messages::init]: Nil answer")
	}
	response := answer.Output["response"].(map[string]interface{})
	lp.key = response["key"].(string)
	lp.server = response["server"].(string)
	lp.ts = int64(response["ts"].(float64))
	lp.NewMessageChan = make(NewMessageChan)

	for {
		lp.Go(chanKit)
	}
}

func (lp *LongPoll) Go(chanKit ChanKit) {
	resp, err := http.Get(fmt.Sprintf(requestBody, lp.server, lp.key, lp.ts, conf.VkTimeout))
	if err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to get response: ", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to read data: ", err)
		return
	}
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to parse data: ", err)
		return
	}
	if response["failed"] != nil {
		lp.Start(chanKit)
		log.Println("[INFO] chanKit has been reinitialized")
	} else {
		type jsonBody struct {
			Failed  int64           `json:"failed"`
			Ts      int64           `json:"ts"`
			Updates [][]interface{} `json:"updates"`
		}
		var body jsonBody
		if err := json.Unmarshal(data, &body); err != nil {
			log.Println("[Error] longPoll::process:", err.Error(), "WebResponse:", string(data))
			return
		}
		for _, update := range body.Updates {
			updateID := update[0].(float64)
			switch updateID {
			//TODO ADD NEW CASES
			case 4: //New message action
				label := update[2].(float64)
				if label == unread || label == unread+chat {
					//If message 1 (Message not read) + 16 (Message sent via chat) = 17
					message := new(Message)
					message.Id = int64(update[1].(float64))
					message.UserId = int64(update[3].(float64))
					message.Text = update[6].(string)
					message.Attachments = make(map[string]interface{})
					for key, value := range update[7].(map[string]interface{}) {
						message.Attachments[key] = value.(string)
					}

					lp.NewMessageChan <- *message
				}
			}
		}
		lp.ts = body.Ts
	}
}
