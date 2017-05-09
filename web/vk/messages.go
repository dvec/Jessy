package vk

import (
	"log"
	"time"
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
	Date        int64
	Out         int64
	UserId      int64
	Text        string
	Attachments map[string]interface{}
}

func GetNewMessages(chanKit ChanKit, messageChan chan<- Message) {
	for {
		chanKit.RequestChan <- Request{"messages.get", map[string]string{
			"time_offset": "1",
		}}
		answer := <- chanKit.AnswerChan
		if answer.Error != nil {
			log.Println("[ERROR] [Messages::GetNewMessages]:", answer.Error)
		}
		if answer.Output["response"] == nil {
			log.Println("[ERROR]	[Messages::GetNewMessages]: Nil answer")
			continue
		}
		response := answer.Output["response"].([]interface{})
		for _, currentMessage := range response[1:] {
			parsedMessage := currentMessage.(map[string]interface{})
			if parsedMessage["read_state"].(float64) == 0 {
				log.Println("[INFO] Got new message:", parsedMessage)
				var attachments map[string]interface{}
				if parsedMessage["attachments"] != nil {
					attachments = parsedMessage["attachments"].(map[string]interface{})
				}
				messageChan <- Message{
					Date:        int64(parsedMessage["date"].(float64)),
					Out:         int64(parsedMessage["out"].(float64)),
					Id:          int64(parsedMessage["mid"].(float64)),
					UserId:      int64(parsedMessage["uid"].(float64)),
					Text:        parsedMessage["body"].(string),
					Attachments: attachments,
				}
			}
		}
		time.Sleep(time.Second / 3)
	}
}