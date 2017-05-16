package commands

import (
	"time"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"fmt"
	"main/engine/commands/tools"
)


func GetState(args tools.FuncArgs) {
	start := time.Now().UnixNano()
	sort.Ints(rand.Perm(1000))
	metering := strconv.FormatInt(time.Now().UnixNano() - start, 10)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": fmt.Sprintf("Я отсортировала массив из 1000 элементов за %v наносекунд", metering),
	})
}

func GetGen(args tools.FuncArgs) {
	var message string
	information := strconv.FormatInt(int64(tools.GetRandomNum(args.Message.Text)%100), 10)
	message = fmt.Sprintf("С вероятностью %v%%", information)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func GetHelp(args tools.FuncArgs) {
	var message string
	if len(args.Message.Text) != 1 {
		message = tools.GetHelp(args.Message.Text, args.DataCache.CommandDataCache.Help)
	} else {
		message = tools.GetHelp("", args.DataCache.CommandDataCache.Help)
	}
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func Bash(args tools.FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Bash.ChooseRandom(),
	})
}

func IThappens(args tools.FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.IThappens.ChooseRandom(),
	})
}

func Zadolbali(args tools.FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Zadolbali.ChooseRandom(),
	})
}

func News(args tools.FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var message string
	if tools.CheckData(words, []string{"s", "i"}) {
		count, _ := strconv.ParseInt(words[1], 10, 8)
		if count > 7 {
			message = "Я не помню столько новостей"
		} else {
			message = strings.Join(args.DataCache.RSSCache.News.Data[:count], "\n")
		}
	} else {
		message = strings.Join(args.DataCache.RSSCache.News.Data[:3], "\n")
	}

	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func Cities(args tools.FuncArgs) {
	args.InterceptIndications.Add(args.Message.UserId)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": "Добро пожаловать в игру 'города'! Начинай!",
	})

	already := []string{}
	var expectedSymbol rune
	for {
		var answer string
		message := <- args.InterceptIndications.InterceptedMessage[args.Message.UserId]
		message.Text = strings.TrimRight(strings.ToLower(message.Text), "ьъый")

		runes := []rune(message.Text)

		args.DataCache.CommandDataCache.Cities.Lock()

		if message.Text == "хватит" {
			args.InterceptIndications.Delete(args.Message.UserId)
			args.ApiChan.MakeRequest("messages.send", map[string]string{
				"user_id": strconv.FormatInt(args.Message.UserId, 10),
				"message": "Игра прекращена. Ты можешь продолжать со мной общение",
			})
			return
		}

		if expectedSymbol != 0 && expectedSymbol != runes[0] {
			answer = "Ты должен назвать слово на букву " + string(expectedSymbol)
			goto SEND
		}

		if tools.Contains(already, message.Text) {
			answer = "Уже было!"
			goto SEND
		}
		already = append(already, message.Text)

		for _, char := range message.Text {
			if int32('а') >  char || char > int32('я') {
				answer = "Используй только русские буквы!"
				goto SEND
			}
		}

		{
			lastSymbol := runes[len(runes) - 1]
			for _, city := range args.DataCache.CommandDataCache.Cities.Data.CitiesList {
				if city.Name == "" {
					continue
				}

				if []rune(strings.ToLower(city.Name))[0] == lastSymbol {
					if !tools.Contains(already, strings.ToLower(city.Name)) {
						answer = city.Name
						answerRunes := []rune(answer)
						expectedSymbol = answerRunes[len(answerRunes) - 1]
						already = append(already, strings.ToLower(answer))
						goto SEND
					}
				}
			}
		}

		answer = "Ты выиграл. Мои поздравления! Я начинаю новую игру. Для выхода напиши 'хватит'"

		SEND: args.ApiChan.MakeRequest("messages.send", map[string]string{
			"user_id": strconv.FormatInt(args.Message.UserId, 10),
			"message": answer,
		})
		args.DataCache.CommandDataCache.Cities.Unlock()
	}
}