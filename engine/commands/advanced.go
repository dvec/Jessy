package commands

import (
	"strconv"
	"fmt"
)


func toSimpleText(text string, banSymbols []string) string {
	var isReady bool
	for !isReady && len(text) != 0 {
		isReady = true
		for _, char := range banSymbols {
			if string(text[len(text) - 1]) == char {
				text = text[:len(text)-1]
				isReady = false
				break
			}
		}
	}
	return text
}

func getRandomNum(text string) int {
	out := 50
	for _, char := range text {
		out += int(char)
	}
	return out % 100
}

func checkData(args []string, filter []string) bool {
	if len(args) != len(filter) {
		return false
	}

	l: for index, word := range args {
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
		return fmt.Sprintf("Справка по команде %v: %v", name, cache[name])
	} else {
		return "Нет справки по такой команде"
	}
}