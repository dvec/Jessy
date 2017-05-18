package main

import ("os"
	"log"
	"time"
	"fmt"
	"main/conf"
	"main/web/vk"
	"main/web/rss"
	"main/kernel"
	"main/kernel/cache"
	"main/kernel/interception"
	"main/kernel/commands"
)

//Is log will be written
const isLogFileWritten = false

var (
	//New log file name
	logFilePath = conf.LOG_DIR_PATH + "/" + fmt.Sprintf("%d", time.Now().Unix()) + ".log"

	//Mandatory files list
	pathWays = []struct{
		path string
		isDir bool
	}{
		{conf.DATA_DIR_PATH, true}, // data
		{conf.DATA_DIR_PATH + "/dict.aiml.xml", false}, // data/dict.aiml.xml
		{conf.LOG_DIR_PATH, true}, // data/log
		{logFilePath, false}, // data/log/xxxxxxxxxx.log
		{conf.COMMANDS_DIR_PATH, true}, // data/commands
		{conf.COMMANDS_DIR_PATH + "/help.xml", false}, // data/commands/help.xml
		{conf.COMMANDS_DIR_PATH + "/cities.xml", false}, // data/commands/cities.xml
	}
)

func main() {
	log.Println("[INFO] main.go started")
	//Checking files
	for _, file := range pathWays {
		if _, err := os.Stat(file.path); os.IsNotExist(err) {
			if file.isDir {
				os.Mkdir(file.path, conf.DATA_FILE_PERMISSION)
			} else {
				os.OpenFile(file.path, os.O_CREATE, conf.DATA_FILE_PERMISSION)
			}
			log.Print("[INFO] ", file.path, " has been created")
		}
	}

	logFile, fileOpenError := os.OpenFile(logFilePath, os.O_RDWR, conf.DATA_FILE_PERMISSION)
	if fileOpenError != nil {
		log.Print("[ERROR] [main::main()] Failed to open log file: ", fileOpenError)
	}

	//Redirecting log output
	//noinspection ALL
	if isLogFileWritten {
		log.Println("[INFO] Output will be redirected to a log file.")
		log.SetOutput(logFile)
	}

	//vk API initialization
	log.Println("[INFO] Initializing vk api...")
	var api vk.Api
	api.AccessToken = conf.VK_TOKEN
	var dataCache cache.DataCache

	//Chan initialization
	messageChan := make(chan vk.Message)
	api.InitChanKit()

	//RSS file (news, bash, etc) initialization
	log.Println("[INFO] Initializing cache...")
	dataCache.InitCache()
	rss.UpdateRss(&dataCache.RssCache)

	//Intercept indications init
	indications := interception.Indications{}
	indications.Init()

	//Runs long poll checking
	var lp vk.LongPoll
	go func() {
		lp.Init(api.ChanKit)
		for {
			lp.Go(api.ChanKit, messageChan)
		}
	}()

	//Chan checking
	for {
		select {
		case message := <- messageChan: //New message
			log.Println("[INFO] New message detected: ", message)
			go kernel.Perform(commands.FuncArgs{api.ChanKit, message, dataCache, indications})
		case request := <- api.ChanKit.RequestChan: //New api request
			out, err := api.Request(request.Name, request.Params)
			api.ChanKit.AnswerChan <- vk.Answer{out, err}
			time.Sleep(time.Second / conf.MAX_REQUEST_PER_SECOND)
		case <- time.After(conf.RSS_UPDATE_DELAY): //Time to update cache
			log.Println("[INFO] Time to update RSS files")
			go rss.UpdateRss(&dataCache.RssCache)
		}
	}
}