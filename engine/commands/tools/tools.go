package tools

import (
	"strconv"
	"fmt"
	"main/engine/cache"
)

func GetRandomNum(text string) int {
	out := 50
	for _, char := range text {
		out += int(char)
	}
	return out
}

func CheckData(args []string, template []string) bool {
	if len(args) < len(template) {
		return false
	}

	l: for index, word := range args {
		if index < len(template) {
			switch template[index] {
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
		var emojiDict = map[string]string{
			"ready": "&#128215;",
			"test": "&#128217;",
			"error": "&#128213;",
			"dev": "&#128216;",

		}
		cache.Lock()
		defer cache.Unlock()
		for _, command := range cache.Data.HelpList {
			commandList += fmt.Sprintf("%v | %v \n", emojiDict[command.State], command.Name)
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

func Contains(arr []string, value string) bool {
	for _, a := range arr {
		if a == value {
			return true
		}
	}
	return false
}