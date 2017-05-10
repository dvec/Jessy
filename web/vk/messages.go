package vk

import (
	"log"
	"net/http"
	"fmt"
	"main/conf"
	"io/ioutil"
	"encoding/json"
)

type Request struct {
	Name	string
	Params 	map[string]string
}

type Answer struct {
	Output 	map[string]interface{}
	Error	error
}

type ChanKit struct {
	RequestChan	chan Request
	AnswerChan	chan Answer
}

func (chanKit ChanKit)MakeRequest(name string, params map[string]string) Answer {
	chanKit.RequestChan <- Request{name, params}
	return <- chanKit.AnswerChan
}

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

func (lp *LongPoll) Go(chanKit ChanKit, messageChan chan<- Message) {
	resp, err := http.Get(fmt.Sprintf("https://%v?act=a_check&key=%v&ts=%v&wait=%v&mode=2&version=1", lp.server, lp.key, lp.ts, conf.TIMEOUT))
	if err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to get response: ", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to read data: ", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("[ERROR] [Messages::Go]: failed to parse data: ", err)
	}
	if response["failed"] != nil {
		log.Println("[INFO] Reinitializating chanKit...")
		lp.Init(chanKit)
	} else {
		type jsonBody struct {
			Failed  int64           `json:"failed"`
			Ts      int64           `json:"ts"`
			Updates [][]interface{} `json:"updates"`
		}
		var body jsonBody
		if err := json.Unmarshal(data, &body); err != nil {
			log.Println("[Error] longPoll::process:", err.Error(), "WebResponse:", string(data))
		}
		for _, update := range body.Updates {
			updateID := update[0].(float64)
			switch updateID {
			//TODO ADD NEW CASES
			case 4: //New message action
				label := update[2].(float64)
				if label == UNREAD || label == UNREAD + CHAT {
					//If message from user 1 (Message not read) + 16 (Message sent via chat) = 17
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
	resp.Body.Close() //I can't move it in defer because this function never ends. Sorry me for bad code. Hold on. We are with you
	lp.Go(chanKit, messageChan)
}