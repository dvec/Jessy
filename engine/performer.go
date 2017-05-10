package engine

import (
	"main/web/vk"
	"strconv"
	"strings"
	"main/engine/commands"
	"log"
	"main/engine/cache"
	"math/rand"
)

type function struct {
	name 		string
	function 	func(commands.FuncArgs)
}

func getAnswer(data map[string][]string, message string) string {
	answer := data[message]
	if len(answer) != 0 {
		return answer[rand.Intn(len(answer))]
	} else {
		words := strings.Split(message, " ")
		var sep string
		if len(words) > 1 {
			sep = " "
		} else {
			sep = ""
		}
		newMessage := strings.Join(words[:len(words) - 1], sep)
		newData := map[string][]string{}
		for request, answers := range data {
			messageLen := len(newMessage)
			if messageLen <= len(request) {
				newData[request[:len(newMessage)]] = answers
			}
		}
		return getAnswer(newData, newMessage)
	}
}

func getFunctions() []function {
	commandsList := []function {
		{"статус", commands.GetState},
		{"помощь", commands.GetHelp},
		{"инфа", commands.GetGen},
		{"баш", commands.Bash},
		{"айти", commands.IThappens},
		{"задолбали", commands.Zadolbali},
		{"новости", commands.News},
	}
	return commandsList
}

func Perform(chanKit vk.ChanKit, message vk.Message, dataCache cache.DataCache) {
	text := strings.ToLower(strings.Trim(message.Text, "?!():.,|"))
	args := strings.Split(text, " ")
	firstWord := strings.ToLower(args[0])
	for _, command := range getFunctions() {
		if firstWord == command.name {
			log.Println("[INFO] Command detected: ", firstWord)
			params := commands.FuncArgs{ApiChan: chanKit, Message: message, DataCache: dataCache}
			command.function(params)
			return
		}
	}
	log.Println("[INFO] No command detected. Running reiteration")
	dataCache.DictionaryCache.Lock()
	answer := getAnswer(dataCache.DictionaryCache.Data, text)
	dataCache.DictionaryCache.Unlock()
	params := map[string]string{
		"user_id": strconv.FormatInt(message.UserId, 10),
	}
	attach := strings.Index(answer, "<attach>")
	if attach != -1 {
		attachEnd := strings.Index(answer, "<end>")
		if attachEnd != -1 {
			params["attachment"] = answer[attach + len("<attach>"):attachEnd]
			params["messsage"] = answer[:attach]
		} else {
			params["message"] = "Internal error"
		}
	} else {
		params["message"] = answer
	}
	chanKit.MakeRequest("messages.send", params)
}