package commands

import (
	"time"
	"math/rand"
	"sort"
	"strconv"
	"main/web/vk"
)

type FuncArgs struct {
	ApiChan vk.ChanKit
	Message vk.Message
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

func Print(args FuncArgs) {
	args.ApiChan.MakeRequest("messages.send", map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": toSimpleText(args.Message.Text, []string{"\\", "|", "/"}),
	})
}