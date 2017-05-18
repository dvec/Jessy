package tools

import (
	"strconv"
	"fmt"
	"main/kernel/cache/cachetypes"
)

const (
	//GET_RANDOM
	startNum	= 50

	//CHECK_DATA
	IntegerTag  	= "i"
	StringTag   	= "s"
	AnythingTag 	= "*"

	//GET_HELP
	DefaultTag      = ""
	ponyDescription = "%v | %v \n"
	simpleAnswer    = "Список моих команд: \n%v" +
		" Вы можете посмотреть справку по любой из них, набрав:" +
		" \nпомощь [название команды]"
	complexAnswer	= `Справка по команде "%v": %v %v`
	sampleIntroText	= "Примеры использования: \n"
	sampleTemplate	= "USER> %v\nJESSY> %v"
	noHelpError	= "Нет справки по такой команде"
)

var emojiDict = map[string]string{
	"ready": "&#128215;",
	"test": "&#128217;",
	"error": "&#128213;",
	"dev": "&#128216;",

}

//Function returns a random value based on a string (text).
func GetRandomNum(text string) int {
	out := startNum
	for _, char := range text {
		out += int(char)
	}
	return out
}

//Function checks if the array of strings (args) matches template (template).
func IfMatch(args []string, template []string) bool {
	if len(args) < len(template) {
		return false
	}

	l: for index, word := range args {
		if index < len(template) {
			switch template[index] {
			case IntegerTag:
				_, err := strconv.ParseInt(word, 10, 64)
				if err != nil {
					return false
				}
			case StringTag:
				//TODO ADD
			case AnythingTag:
				break l
			}
		}

		if index > len(template) {
			return false
		}
	}

	return true
}

//Function returns function reference with a name (name). Also you must pass HelpCache to boost speed.
func GetHelp(name string, cache cachetypes.HelpCache) string {
	if name == DefaultTag {
		var commandList string
		cache.Lock()
		defer cache.Unlock()
		for _, command := range cache.Data.HelpList {
			commandList += fmt.Sprintf(ponyDescription, emojiDict[command.State], command.Name)
		}
		return fmt.Sprintf(simpleAnswer, commandList)
	}
	cache.Lock()
	defer cache.Unlock()
	for _, command := range cache.Data.HelpList {
		if command.Name == name {
			samples := ""
			if command.Samples != nil {
				samples = sampleIntroText
				for _, sample := range command.Samples {
					samples += fmt.Sprintf(sampleTemplate, sample.Body, sample.Out)
				}
			}
			return fmt.Sprintf(complexAnswer, name, command.Description, samples)
		}
	}
	return noHelpError
}

//Function returns true if array (arr) contains value (value) and false in other cases
func Contains(arr []string, value string) bool {
	for _, a := range arr {
		if a == value {
			return true
		}
	}
	return false
}