package commands

import (
	"time"
	"math/rand"
	"sort"
	"strconv"
	"main/web/vk"
	"strings"
	"main/engine/cache"
)

type FuncArgs struct {
	ApiChan vk.ChanKit
	Message vk.Message
	DataCache cache.DataCache
}


func GetState(args FuncArgs) {
	start := time.Now().UnixNano()
	sort.Ints(rand.Perm(1000))
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": strconv.FormatInt(time.Now().UnixNano() - start, 10),
	})
}

func GetGen(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": strconv.FormatInt(int64(getRandomNum(args.Message.Text)), 10),
	})
}

func GetHelp(args FuncArgs) {
	words := strings.Split(args.Message.Text, " ")
	var name string
	if checkData(words, []string{"s", "s"}) {
		name = words[1]
	} else {
		name = ""
	}
	message := getHelp(name, args.DataCache.CommandDataCache.Help.Data)
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
	})
}

func Print(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": toSimpleText(args.Message.Text, []string{"\\", "|", "/"}),
	})
}

func Bash(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": args.DataCache.RSSCache.Bash.Data[rand.Intn(len(args.DataCache.RSSCache.Bash.Data) - 1)],
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