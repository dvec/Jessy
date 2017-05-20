package functions

import (
	"main/web/vk"
	"main/kernel/cache"
	"main/kernel/interception"
	"strconv"
	"main/conf"
	"errors"
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
func (args *FuncArgs) Reply(message string, attach ...string) error {
	if len(message) >= conf.MaxMessageLen {
		return errors.New("Length too long")
	}
	var attachments string

	if len(attach) != 0 {
		attachments = attach[0]
	}

	args.ApiChan.MakeRequest(messageSendMethod, map[string]string{
		"user_id":     strconv.FormatInt(args.Message.UserId, 10),
		"message":     message,
		"attachment": attachments,
	})

	return nil
}