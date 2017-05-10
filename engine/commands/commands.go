package commands

import (
	"time"
	"math/rand"
	"sort"
	"strconv"
	"main/web/vk"
	"strings"
	"main/engine/cache"
	"fmt"
)

type FuncArgs struct {
	ApiChan vk.ChanKit
	Message vk.Message
	DataCache cache.DataCache
}


func GetState(args FuncArgs) {
	start := time.Now().UnixNano()
	sort.Ints(rand.Perm(1000))
	metering := strconv.FormatInt(time.Now().UnixNano() - start, 10)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": fmt.Sprintf("Я отсортировала массив из 1000 элементов за %v наносекунд", metering),
	})
}

func GetGen(args FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var message string
	if checkData(words, []string{"s", "*"}) {
		information := strconv.FormatInt(int64(getRandomNum(args.Message.Text)%100), 10)
		message = fmt.Sprintf("С вероятностью %v%%", information)
	} else {
		message = getHelp("инфа", args.DataCache.CommandDataCache.Help)
	}
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func GetHelp(args FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var name string
	if checkData(words, []string{"s", "s"}) {
		name = strings.Join(words[1:], " ")
	} else {
		name = ""
	}
	message := getHelp(name, args.DataCache.CommandDataCache.Help)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func Bash(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Bash.Data[rand.Intn(len(args.DataCache.RSSCache.Bash.Data) - 1)],
	})
}

func IThappens(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.IThappens.Data[rand.Intn(len(args.DataCache.RSSCache.IThappens.Data) - 1)],
	})
}

func Zadolbali(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Zadolbali.Data[rand.Intn(len(args.DataCache.RSSCache.Zadolbali.Data) - 1)],
	})
}

func News(args FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var message string
	if checkData(words, []string{"s", "i"}) {
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