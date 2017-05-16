package engine

import (
	"strconv"
	"strings"
	"main/engine/commands"
	"log"
	"fmt"
	"main/engine/commands/tools"
)

var commandsList = map[string]func(args tools.FuncArgs) {
	"state": commands.GetState,
	"help": commands.GetHelp,
	"cities": commands.Cities,
	"information": commands.GetGen,
	"bash": commands.Bash,
	"ithappens": commands.IThappens,
	"zadolbali": commands.Zadolbali,
	"news": commands.News,
}

func findFlags(text string, flags []string) (map[string]string, string) {
	out := make(map[string]string)
	minFlagindex := len(text)
	for _, flag := range flags {
		begin := fmt.Sprintf("<%v>", flag)
		end := fmt.Sprintf("</%v>", flag)
		beginIndex := strings.Index(text, begin)
		endIndex := strings.Index(text, end)
		if beginIndex < endIndex {
			out[flag] = text[beginIndex + len(begin):endIndex]
			if minFlagindex > beginIndex {
				minFlagindex = beginIndex
			}
		}
	}

	return out, text[:minFlagindex]
}

func checkInterceptIndications(args tools.FuncArgs) bool {
	args.InterceptIndications.Lock()
	defer args.InterceptIndications.Unlock()
	if args.InterceptIndications.InterceptedMessage[args.Message.UserId] != nil {
		args.InterceptIndications.InterceptedMessage[args.Message.UserId] <- args.Message
		return true
	}
	return false
}

func Perform(args tools.FuncArgs) {
	if checkInterceptIndications(args) { return }
	text := strings.ToLower(strings.Trim(args.Message.Text, "?!():.,|"))
	log.Println("[INFO] No command detected. Running performation")
	args.DataCache.DictionaryCache.Lock()
	answer, err := args.DataCache.DictionaryCache.Data.Respond(strings.ToLower(text))
	args.DataCache.DictionaryCache.Unlock()
	if err != nil {
		log.Println("[ERROR] [main::engine::performer.go] Failed to get answer: ", err)
	}

	params := map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
	}
	flags, newMessage := findFlags(answer, []string{"attach", "call"})
	params["message"] = newMessage
	params["attachment"] = flags["attach"]
	if flags["call"] != "" {
		funcParams := strings.Split(flags["call"], "|")
		name := funcParams[0]

		if commandsList[name] != nil {
			args.Message.Text = funcParams[1]
			commandsList[name](args)
			return
		} else {
			params["message"] = "Internal error"
		}
	}

	args.ApiChan.MakeRequest("messages.send", params)
}