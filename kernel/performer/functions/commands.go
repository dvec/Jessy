package functions

import (
	"main/web/vk"
	"main/kernel/cache"
	"main/kernel/interception"
	"strconv"
)

const messageSendMethod = "messages.send"

//args for bot commands
type FuncArgs struct {
	ApiChan vk.ChanKit //ChanKit for safety API requests
	Message vk.Message //Message from user
	DataCache cache.DataCache //Cache of data
	InterceptIndications interception.Indications //Interception tools
}

//Send message back to the user
func (args *FuncArgs) Reply(message string, attach ...string) {
	var attachments string

	if len(attachments) != 0 {
		attachments = attach[0]
	}

	args.ApiChan.MakeRequest(messageSendMethod, map[string]string{
		"user_id": strconv.FormatInt(args.Message.UserId, 10),
		"message": message,
		"attachments": attachments,
	})
}