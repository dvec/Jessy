package commands

import (
	"time"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"fmt"
	"main/engine/commands/tools"
	"main/web/vk"
	"main/engine/cache"
	"main/engine/commands/interception"
)

type FuncArgs struct {
	ApiChan vk.ChanKit
	Message vk.Message
	DataCache cache.DataCache
	InterceptIndications interception.Indications
}

const (
	//GET_STATE
	getStateAnswer		= "Я отсортировала массив из 1000 элементов за %v наносекунд"

	//GET_GEN
	getGenAnswer		= "С вероятностью %v%%"

	//GET_HELP
	defaultHelpRequest	= ""

	//NEWS
	tooMuchNewsCountError	= "Я не помню столько новостей"

	//CITIES
	endCommand		= "хватит"
	endMessage		= "Игра прекращена. Ты можешь продолжать со мной общение"
	cityCutset		= "ьъый"
	winMessage		= "Ты выиграл. Мои поздравления! Я начинаю новую игру"
	alreadyError		= "Уже было!"
	badInputError 		= "Используй только русские буквы!"
	welcomeMessage		= "Добро пожаловать в игру 'города'! Для выхода напиши 'хватит'. Начинай!"
	badInputLenError	= "Город не может состоять только из этих букв"
	incorrectSymbolError 	= "Ты должен назвать слово на букву %v"
)

func GetState(args FuncArgs) {
	start := time.Now().UnixNano()
	sort.Ints(rand.Perm(1000))
	metering := strconv.FormatInt(time.Now().UnixNano() - start, 10)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": fmt.Sprintf(getStateAnswer, metering),
	})
}

func GetGen(args FuncArgs) {
	var message string
	information := strconv.FormatInt(int64(tools.GetRandomNum(args.Message.Text)%100), 10)
	message = fmt.Sprintf(getGenAnswer, information)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func GetHelp(args FuncArgs) {
	var message string
	if len(args.Message.Text) != 1 {
		message = tools.GetHelp(args.Message.Text, args.DataCache.CommandDataCache.Help)
	} else {
		message = tools.GetHelp(defaultHelpRequest, args.DataCache.CommandDataCache.Help)
	}
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func Bash(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Bash.ChooseRandom(),
	})
}

func IThappens(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.IThappens.ChooseRandom(),
	})
}

func Zadolbali(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Zadolbali.ChooseRandom(),
	})
}

func News(args FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var message string
	if tools.CheckData(words, []string{"s", "i"}) {
		count, _ := strconv.ParseInt(words[1], 10, 8)
		if count > 7 {
			message = tooMuchNewsCountError
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

func Cities(args FuncArgs) {
	args.InterceptIndications.Add(args.Message.UserId)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": welcomeMessage,
	})

	already := []string{}
	var expectedSymbol rune
	for {
		var answer string
		message := <- args.InterceptIndications.InterceptedMessage[args.Message.UserId]
		message.Text = strings.TrimRight(strings.ToLower(message.Text), cityCutset)

		runes := []rune(message.Text)

		args.DataCache.CommandDataCache.Cities.Lock()

		if len(message.Text) == 0 {
			answer = badInputLenError
			goto SEND
		}

		if message.Text == endCommand {
			args.InterceptIndications.Delete(args.Message.UserId)
			args.ApiChan.MakeRequest("messages.send", map[string]string{
				"user_id": strconv.FormatInt(args.Message.UserId, 10),
				"message": endMessage,
			})
			return
		}

		if expectedSymbol != 0 && expectedSymbol != runes[0] {
			answer = fmt.Sprintf(incorrectSymbolError, expectedSymbol)
			goto SEND
		}

		if tools.Contains(already, message.Text) {
			answer = alreadyError
			goto SEND
		}
		already = append(already, message.Text)

		for _, char := range message.Text {
			if int32('а') >  char || char > int32('я') { //checks if text contains only russian symbols
				answer = badInputError
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
						answerRunes := []rune(strings.TrimRight(answer, cityCutset))
						expectedSymbol = answerRunes[len(answerRunes) - 1]
						already = append(already, strings.ToLower(answer))
						goto SEND
					}
				}
			}
		}

		answer = winMessage

		SEND: args.ApiChan.MakeRequest("messages.send", map[string]string{
			"user_id": strconv.FormatInt(args.Message.UserId, 10),
			"message": answer,
		})
		args.DataCache.CommandDataCache.Cities.Unlock()
	}
}