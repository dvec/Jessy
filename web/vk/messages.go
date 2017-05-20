package vk

import (
	"log"
	"net/http"
	"fmt"
	"main/conf"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type Message struct {
	Id          int64 //Message ID (not User ID)
	UserId      int64 //User ID
	Text        string //Message text
	Attachments map[string]interface{} //Attachments
}

//Chan to inform about a new message
type NewMessageChan chan Message

type LongPoll struct {
	key    string //Long Poll server key
	server string //Long Poll server address
	ts     int64  //Last synchronization time
	NewMessageChan //Chan to inform about a new message
}

const (
	//FLAGS
	unread = 1  //Unread message flag
	chat   = 16 //Chat message flag

	//LONG POLL
	requestBody = "https://%v?act=a_check&key=%v&ts=%v&wait=%v&mode=2&version=1" //API url for post request
)

//Make API data reload (or load)
func (lp *LongPoll) reload(chanKit ChanKit) error {
	answer := chanKit.MakeRequest("messages.getLongPollServer", nil)

	if answer.Error != nil {
		return answer.Error
	}

	if answer.Output["response"] == nil { //If API request returns null data
		return errors.New("Nil answer")
	}

	response := answer.Output["response"].(map[string]interface{})
	lp.key = response["key"].(string)
	lp.server = response["server"].(string)
	lp.ts = int64(response["ts"].(float64))

	return nil
}

//Long Poll initialize function
func (lp *LongPoll) init(chanKit ChanKit) error {
	lp.NewMessageChan = make(NewMessageChan)

	//API data loading (reload function does it)
	if err := lp.reload(chanKit); err != nil {
		return err
	}

	return nil
}

//Starts Long Poll checking
func (lp *LongPoll) Start(chanKit ChanKit) {
	//Initializes Long Poll
	lp.init(chanKit)
	for {
		//Runs Long Poll
		if err := lp.Check(chanKit); err != nil {
			log.Println("[ERROR] Failed to check Long Poll: ", err)
		}
	}
}

//Checks for a Long Poll
func (lp *LongPoll) Check(chanKit ChanKit) error {
	//Receive new messages or timeout
	resp, err := http.Get(fmt.Sprintf(requestBody, lp.server, lp.key, lp.ts, conf.VkTimeout))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Response answer reading
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Response answer parsing
	var body struct {
		Failed  int64           `json:"failed"`
		Ts      int64           `json:"ts"`
		Updates [][]interface{} `json:"updates"`
	}
	if err := json.Unmarshal(data, &body); err != nil {
		log.Println("[Error] longPoll::process:", err.Error(), "WebResponse:", string(data))
		return err
	}

	//Sometimes the TS becomes obsolete and the API sends failed message.
	//We need to update the TS parameter if it occur
	if body.Failed != 0 {
		log.Println("[INFO] chanKit has been reinitialized")
		if err := lp.reload(chanKit); err != nil {
			return err
		}
	} else {
		for _, update := range body.Updates {
			//Update sample:
			//[4 5007 17 2.61220573e+08 1.495282546e+09  ...  test map[]]
			updateID := update[0].(float64)
			switch updateID {
			//TODO ADD NEW CASES
			case 4: //New message action
				label := update[2].(float64)
				if label == unread || label == unread+chat {
					//If message 1 (Message not read) + 16 (Message sent via chat) = 17

					//Message parsing and loading into struct
					message := new(Message)
					message.Id = int64(update[1].(float64))
					message.UserId = int64(update[3].(float64))
					message.Text = update[6].(string)
					message.Attachments = make(map[string]interface{})

					//Performing attachments
					for key, value := range update[7].(map[string]interface{}) {
						message.Attachments[key] = value.(string)
					}

					//Sends a message to chan
					lp.NewMessageChan <- *message
				}
			}
		}
		//Updating TS
		lp.ts = body.Ts
	}

	return nil
}
