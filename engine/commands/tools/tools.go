package tools

import (
	"strconv"
	"fmt"
	"main/engine/cache"
	"main/web/vk"
	"main/engine/commands/interception"
)

type FuncArgs struct {
	ApiChan vk.ChanKit
	Message vk.Message
	DataCache cache.DataCache
	InterceptIndications interception.Indications
}

func GetRandomNum(text string) int {
	out := 50
	for _, char := range text {
		out += int(char)
	}
	return out
}

func CheckData(args []string, filter []string) bool {
	if len(args) < len(filter) {
		return false
	}

	l: for index, word := range args {
		if index < len(filter) {
			switch filter[index] {
			case "i":
				_, err := strconv.ParseInt(word, 10, 64)
				if err != nil {
					return false
				}
			case "*":
				break l
			}
		}
	}

	return true
}

func GetHelp(name string, cache cache.HelpCache) string {
	if name == "" {
		var commandList string
		cache.Lock()
		defer cache.Unlock()
		for _, command := range cache.Data.HelpList {
			var emoji string
			switch command.State {
			case "ready":
				emoji = "&#128215;"
			case "test":
				emoji = "&#128217;"
			case "error":
				emoji = "&#128213;"
			case "dev":
				emoji = "&#128216;"
			}
			commandList += fmt.Sprintf("%v | %v \n", emoji, command.Name)
		}
		return fmt.Sprintf("Список моих команд: \n%v" +
			" Вы можете посмотреть справку по любой из них, набрав:" +
			" \nпомощь [название команды]", commandList)
	}
	cache.Lock()
	defer cache.Unlock()
	for _, command := range cache.Data.HelpList {
		if command.Name == name {
			samples := ""
			if command.Samples != nil {
				samples = "Примеры использования: \n"
				for _, sample := range command.Samples {
					samples += fmt.Sprintf("USER> %v\nJESSY> %v", sample.Body, sample.Out)
				}
			}
			return fmt.Sprintf("Справка по команде "+
				"\"%v\": %v %v", name, command.Description, samples)
		}
	}
	return "Нет справки по такой команде"
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}