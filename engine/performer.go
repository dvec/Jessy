package engine

import (
	"main/web/vk"
	"strconv"
	"strings"
	"main/engine/commands"
	"log"
	"main/engine/cache"
)

type function struct {
	name 		string
	function 	func(commands.FuncArgs)
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
	answer, err := dataCache.DictionaryCache.Data.Respond(strings.ToLower(text))
	dataCache.DictionaryCache.Unlock()
	if err != nil {
		log.Println("[ERROR] [main::engine::performer.go] Failed to get answer: ", err)
	}

	params := map[string]string{
		"user_id": strconv.FormatInt(message.UserId, 10),
	}
	attach := strings.Index(answer, "<attach>")
	if attach != -1 {
		attachEnd := strings.Index(answer, "</attach>")
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