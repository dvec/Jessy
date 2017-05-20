package performer

import (
	"strings"
	"main/kernel/performer/functions"
	"log"
	"fmt"
	"main/kernel/commands"
)

const (
	callSep			= "|"
	callFlag		= "call"
	attachFlag 		= "attach"
	inputCutset		= "?!():.,|"
	internalErrorMessage	= "Internal error"
)


var CommandsList = map[string]func(args functions.FuncArgs) {
	"state":       commands.GetState,
	"help":        commands.GetHelp,
	"cities":      commands.Cities,
	"information": commands.GetGen,
	"bash":        commands.Bash,
	"ithappens":   commands.IThappens,
	"zadolbali":   commands.Zadolbali,
	"news":        commands.News,
}

//This function find flags in the text
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

//This function checks if need to doing interception for this message
func checkInterceptIndications(args functions.FuncArgs) bool {
	args.InterceptIndications.Lock()
	defer args.InterceptIndications.Unlock()
	if args.InterceptIndications.InterceptedMessage[args.Message.UserId] != nil {
		return true
	}
	return false
}

//Main performing function
func Perform(args functions.FuncArgs) {
	if checkInterceptIndications(args) {
		args.InterceptIndications.InterceptedMessage[args.Message.UserId] <- args.Message
		return
	}

	text := strings.ToLower(strings.Trim(args.Message.Text, inputCutset))
	args.DataCache.DictionaryCache.Lock()
	answer, err := args.DataCache.DictionaryCache.Data.Respond(strings.ToLower(text)) //Gets answer
	args.DataCache.DictionaryCache.Unlock()
	if err != nil {
		log.Println("[ERROR] [main::kernel::performer.go] Failed to get answer: ", err)
	}

	flags, message := findFlags(answer, []string{attachFlag, callFlag}) //Checks flags
	//Checks command call
	if flags["call"] != "" {
		//Gets function params
		funcParams := strings.Split(flags[callFlag], callSep)
		name := funcParams[0]

		//Runs command
		if CommandsList[name] != nil {
			args.Message.Text = funcParams[1]
			CommandsList[name](args)
			return
		} else {
			message = internalErrorMessage
		}
	}

	args.Reply(message, flags["attach"]) //Sends answer to the user
}