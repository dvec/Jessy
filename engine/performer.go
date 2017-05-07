package engine

import (
	"main/web/vk"
	"strconv"
	"strings"
	"main/engine/commands"
)


type command struct {
	name 		string
	function 	func(commands.FuncArgs)
}

func Perform(chanKit vk.ChanKit, message vk.Message) {
	args := strings.Split(message.Text, " ")

	for _, command := range getCommands() {
		if args[0] == command.name {
			params := commands.FuncArgs{ApiChan: chanKit, Message: message}
			command.function(params)
			return
		}
	}
	chanKit.MakeRequest("messages.send", map[string]string{
	"user_id":	strconv.FormatInt(message.UserId, 10),
	"message":	message.Text,
	})
}

func getCommands() []command {
	commandsList := []command{
		{"статус", commands.GetState},
		{"напиши", commands.Print},
		{"инфа", commands.GetGen},
	}
	return commandsList
}