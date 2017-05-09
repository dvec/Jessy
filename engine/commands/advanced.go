package commands

import (
	"strconv"
	"fmt"
)


func getRandomNum(text string) int {
	out := 50
	for _, char := range text {
		out += int(char)
	}
	return out
}

func checkData(args []string, filter []string) bool {
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

func getHelp(name string, cache map[string]string) string {
	if name == "" {
		var commandList string
		for command := range cache {
			commandList += fmt.Sprintf("&#128217;|%v \n", command)
		}
		return fmt.Sprintf("Список моих команд: \n%v" +
			" Вы можете посмотреть справку по любой из них, набрав:" +
			" \nпомощь [название команды]", commandList)
	}
	if cache[name] != "" {
		return fmt.Sprintf("Справка по команде: \n" +
			"\"%v\": %v", name, cache[name])
	} else {
		fmt.Println(name)
		return "Нет справки по такой команде"
	}
}