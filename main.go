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
	logFilePath = fmt.Sprintf("%v/%v.log", conf.LogDirPath, strconv.FormatInt(time.Now().Unix(), 10))

	//Mandatory files list
	pathWays = []struct{
		path string
		isDir bool
	}{
		{conf.DataDirPath, true},                      // data
		{conf.DataDirPath + "/dict.aiml.xml", false},  // data/dict.aiml.xml
		{conf.LogDirPath, true},                       // data/log
		{logFilePath, false},                          // data/log/xxxxxxxxxx.log
		{conf.CommandsDirPath, true},                  // data/commands
		{conf.CommandsDirPath + "/help.xml", false},   // data/commands/help.xml
		{conf.CommandsDirPath + "/cities.xml", false}, // data/commands/cities.xml
	}
)

func main() {
	log.Println("[INFO] main.go started")
	//Checking files
	for _, path := range pathWays {
		if _, err := os.Stat(path.path); os.IsNotExist(err) {
			if path.isDir { //Create dir if not exist
				os.Mkdir(path.path, conf.DataFilePermission)
			} else { //Create path if not exist
				os.OpenFile(path.path, os.O_CREATE, conf.DataFilePermission)
			}
			log.Print("[INFO] ", path.path, " has been created") //Log path/dir creation
		}
	}

	logFile, fileOpenError := os.OpenFile(logFilePath, os.O_RDWR, conf.DataFilePermission)
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
	api.AccessToken = conf.VkToken
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
		case request := <- api.ChanKit.RequestChan:                   //New api request
			out, err := api.Request(request.Name, request.Params) //Request API method
			api.ChanKit.AnswerChan <- vk.Answer{out, err}         //Sending API answer back
			time.Sleep(time.Second / conf.MaxRequestPerSecond)    //Delay
		case <- time.After(conf.RssUpdateDelay): //Time to update cache
			log.Println("[INFO] Time to update RSS files")
			go rss.UpdateRss(&dataCache.RssCache)  //Updating RSS in new thread
		}
	}
}