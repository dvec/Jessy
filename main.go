package main

import ("os"
	"log"
	"time"
	"fmt"
	"main/conf"
	"main/web/vk"
	"main/web/rss"
	"main/engine"
	"main/engine/cache"
)

func main() {
	log.Println("[INFO] main.go started")
	logFilePath := conf.LOG_DIR_PATH + "/" + fmt.Sprintf("%d", time.Now().Unix()) + ".log"
	pathWays := []struct{
		path string
		isDir bool
	}{
		{conf.DATA_DIR_PATH, true},
		{conf.DATA_DIR_PATH + "/dictionary.bin", false},
		{conf.RSS_DIR_PATH, true},
		{conf.RSS_DIR_PATH + "/bash.dat", false},
		{conf.RSS_DIR_PATH + "/news.dat", false},
		{conf.RSS_DIR_PATH + "/ithappens.dat", false},
		{conf.RSS_DIR_PATH + "/zadolbali.dat", false},
		{conf.LOG_DIR_PATH, true},
		{logFilePath, false},
		{conf.COMMANDS_DIR_PATH, true},
		{conf.COMMANDS_DIR_PATH + "/help.dat", false},
	}

	log.Println("[INFO] File configuring has been finished. Starting check files.")
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
	log.Println("[INFO] Updating RSS files...")
	rss.Update()

	log.Println("[INFO] Initializating vk api...")
	var api vk.Api
	api.AccessToken = conf.TOKEN
	var dataCache cache.DataCache

	messageChan := make(chan vk.Message)
	api.InitChanKit()

	log.Println("[INFO] Initializating cache...")
	dataCache.InitCache()

	log.Print("[INFO] Initialization is complete. Output will be redirected to a log file")
	log.SetOutput(logFile)
	go vk.GetNewMessages(api.ChanKit, messageChan)
	go func() {
		for {
			select {
			case message := <- messageChan:
				go engine.Perform(api.ChanKit, message, dataCache)
			case request := <- api.ChanKit.RequestChan:
				out, err := api.Request(request.Name, request.Params)
				api.ChanKit.AnswerChan <- vk.Answer{out, err}
				time.Sleep(time.Second / 3)
			case <- time.After(time.Minute * 5):
				log.Println("[INFO] Time to update RSS files")
				go rss.Update()
				//noinspection GoDeferInLoop
				defer func() {
					go dataCache.UpdateCache()
				}()
			}
		}
	}()
	fmt.Scanln()
}