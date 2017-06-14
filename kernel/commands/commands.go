package commands

import (
	"strings"
	"main/kernel/commands/tools"
	"strconv"
	"sort"
	"time"
	"fmt"
	"math/rand"
	"main/kernel/performer/functions"
)

const (
	//GET_STATE
	getStateAnswer = "Я отсортировала массив из 1000 элементов за %v наносекунд"

	//GET_GEN
	getGenAnswer = "С вероятностью %v%%"

	//NEWS
	tooMuchNewsCountError = "Я не помню столько новостей"
	badNewsCountError = "Я не могу сказать тебе столько новостей. Сам попробуй!"

	//CITIES
	endCommand = "хватит"
	endMessage = "Игра прекращена. Ты можешь продолжать со мной общение"
	cityCutset = "ьъый" //To delete these strings from the city name
	winMessage = "Ты выиграл. Мои поздравления! Я начинаю новую игру"
	alreadyError = "Уже было!"
	welcomeMessage = `Добро пожаловать в игру 'города'! Для выхода напиши "хватит". Начинай!`
	badInputLenError = "Город не может состоять только из этих букв"
	incorrectSymbolError = `Ты должен назвать слово на букву "%v"`
)

func GetState(args functions.FuncArgs) {
	start := time.Now().UnixNano()
	sort.Ints(rand.Perm(1000))
	metering := strconv.FormatInt(time.Now().UnixNano() - start, 10)
	args.Reply(fmt.Sprintf(getStateAnswer, metering))
}

func GetGen(args functions.FuncArgs) {
	var message string
	information := strconv.FormatInt(int64(tools.GetRandomNum(args.Message.Text)%100), 10)
	message = fmt.Sprintf(getGenAnswer, information)
	args.Reply(message)
}

func GetHelp(args functions.FuncArgs) {
	var message string
	if len(args.Message.Text) != 1 {
		message = tools.GetHelp(args.Message.Text, args.DataCache.CommandDataCache.Help)
	} else {
		message = tools.GetHelp(tools.DefaultTag, args.DataCache.CommandDataCache.Help)
	}
	args.Reply(message)
}

func Bash(args functions.FuncArgs) {
	args.Reply(args.DataCache.RssCache.Bash.ChooseRandom())
}

func IThappens(args functions.FuncArgs) {
	args.Reply(args.DataCache.RssCache.IThappens.ChooseRandom())
}

func Zadolbali(args functions.FuncArgs) {
	args.Reply(args.DataCache.RssCache.Zadolbali.ChooseRandom())
}

func News(args functions.FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var message string
	if tools.IfMatch(words, []string{"i"}) {
		count, _ := strconv.ParseInt(words[0], 10, 8)
		if count > 7 {
			message = tooMuchNewsCountError
		} else if count <= 0 {
			message = badNewsCountError
		} else {
			message = strings.Join(args.DataCache.RssCache.News.Data[:count], "\n")
		}
	} else {
		message = strings.Join(args.DataCache.RssCache.News.Data[:3], "\n")
	}

	args.Reply(message)
}

func Cities(args functions.FuncArgs) {
	args.InterceptIndications.Add(args.Message.UserId)
	args.Reply(welcomeMessage)

	already := []string{}
	var expectedSymbol rune
	for {
		func() {
			var answer string
			message := <-args.InterceptIndications.InterceptedMessage[args.Message.UserId]
			message.Text = strings.TrimRight(strings.ToLower(message.Text), cityCutset)

			runes := []rune(message.Text)

			if len(message.Text) == 0 {
				args.Reply(badInputLenError)
				return
			}

			if message.Text == endCommand {
				args.InterceptIndications.Delete(args.Message.UserId)
				args.Reply(endMessage)
				return
			}

			if expectedSymbol != 0 && expectedSymbol != runes[0] {
				args.Reply(fmt.Sprintf(incorrectSymbolError, string(expectedSymbol)))
				return
			}

			if tools.Contains(already, message.Text) {
				args.Reply(alreadyError)
				return
			}
			already = append(already, message.Text)

			args.DataCache.CommandDataCache.Cities.Lock()
			defer args.DataCache.CommandDataCache.Cities.Unlock()
			lastSymbol := runes[len(runes)-1]
			for _, city := range args.DataCache.CommandDataCache.Cities.Data.CitiesList {
				if city.Name == "" {
					args.DataCache.CommandDataCache.Cities.Unlock()
					return
				}

				answerRunes := []rune(strings.TrimRight(strings.ToLower(city.Name), cityCutset))
				fmt.Println(string(answerRunes[0]), string(lastSymbol))
				if answerRunes[0] == lastSymbol {
					if !tools.Contains(already, strings.ToLower(city.Name)) {
						expectedSymbol = answerRunes[len(answerRunes)-1]
						already = append(already, strings.ToLower(answer))
						args.Reply(city.Name)
						return
					}
				}
			}

			already = []string{}
			args.Reply(winMessage)
		}()
	}
}