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

type LongPoll struct {
	key	string
	server	string
	ts	int64
}

type NewMessageChan chan<- Message

const (
	UNREAD = 1
	CHAT = 16
)

func (lp *LongPoll) Init(chanKit ChanKit) {
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
}

func (lp *LongPoll) Go(chanKit ChanKit, messageChan NewMessageChan) {
	resp, err := http.Get(fmt.Sprintf("https://%v?act=a_check&key=%v&ts=%v&wait=%v&mode=2&version=1", lp.server, lp.key, lp.ts, conf.VK_TIMEOUT))
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
		lp.Init(chanKit)
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
				if label == UNREAD || label == UNREAD + CHAT {
					//If message 1 (Message not read) + 16 (Message sent via chat) = 17
					message := new(Message)
					message.Id = int64(update[1].(float64))
					message.UserId = int64(update[3].(float64))
					message.Text = update[6].(string)
					message.Attachments = make(map[string]interface{})
					for key, value := range update[7].(map[string]interface{}) {
						message.Attachments[key] = value.(string)
					}
					messageChan <- *message
				}
			}
		}
		lp.ts = body.Ts
	}
}

func (lp *LongPoll) Start(kit ChanKit, messageChan NewMessageChan) {
	for {
		lp.Go(kit, messageChan)
	}
}
