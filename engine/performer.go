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

func Perform(chanKit vk.ChanKit, message vk.Message, dataCache cache.DataCache) {
	args := strings.Split(message.Text, " ")
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
	chanKit.MakeRequest("messages.send", map[string]string{
		"user_id":	strconv.FormatInt(message.UserId, 10),
		"message":	message.Text,
	})
}

func getFunctions() []function {
	commandsList := []function {
		{"статус", commands.GetState},
		{"помощь", commands.GetHelp},
		{"напиши", commands.Print},
		{"инфа", commands.GetGen},
		{"баш", commands.Bash},
		{"новости", commands.News},
	}
	return commandsList
}