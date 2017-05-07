package commands

import (
	"strconv"
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

func parseData(args []string, filter []string) bool {
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