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
	"strconv"
)

//Is log will be written
const isLogFileWritten = false

var (
	//New log file name
	logFilePath = fmt.Sprintf("%v/%v.log", conf.LOG_DIR_PATH, strconv.FormatInt(time.Now().Unix(), 10))

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
	for _, path := range pathWays {
		if _, err := os.Stat(path.path); os.IsNotExist(err) {
			if path.isDir { //Create dir if not exist
				os.Mkdir(path.path, conf.DATA_FILE_PERMISSION)
			} else { //Create path if not exist
				os.OpenFile(path.path, os.O_CREATE, conf.DATA_FILE_PERMISSION)
			}
			log.Print("[INFO] ", path.path, " has been created") //Log path/dir creation
		}
	}

	logFile, fileOpenError := os.OpenFile(logFilePath, os.O_RDWR, conf.DATA_FILE_PERMISSION)
	if fileOpenError != nil {
		log.Print("[ERROR] [main::main()] Failed to open log path: ", fileOpenError)
	}

	//noinspection ALL
	if isLogFileWritten {
		log.Println("[INFO] Output will be redirected to a log path.")
		log.SetOutput(logFile) //Redirecting log output
	}

	//vk API initialization
	log.Println("[INFO] Initializing vk api...")
	var api vk.Api
	api.AccessToken = conf.VK_TOKEN
	var dataCache cache.DataCache

	//Chan initialization
	log.Println("[INFO] Initializating chan kit...")
	api.InitChanKit()

	//RSS path (news, bash, etc) initialization
	log.Println("[INFO] Initializing cache...")
	dataCache.InitCache()
	rss.UpdateRss(&dataCache.RssCache)

	//Intercept indications init
	indications := interception.Indications{}
	indications.Init()

	//Runs long poll checking
	var lp vk.LongPoll
	go lp.Start(api.ChanKit)

	//Chan checking
	for {
		select {
		case message := <- lp.NewMessageChan: //New message
			log.Println("[INFO] New message detected: ", message)
			args := commands.FuncArgs{
				ApiChan: api.ChanKit, Message: message,
				DataCache: dataCache, InterceptIndications: indications,
			} //Creating func params
			go kernel.Perform(args)
		case request := <- api.ChanKit.RequestChan: //New api request
			out, err := api.Request(request.Name, request.Params) //Request API method
			api.ChanKit.AnswerChan <- vk.Answer{out, err} //Sending API answer back
			time.Sleep(time.Second / conf.MAX_REQUEST_PER_SECOND) //Delay
		case <- time.After(conf.RSS_UPDATE_DELAY): //Time to update cache
			log.Println("[INFO] Time to update RSS files")
			go rss.UpdateRss(&dataCache.RssCache)  //Updating RSS in new thread
		}
	}
}