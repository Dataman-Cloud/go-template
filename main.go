package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Dataman-Cloud/go-template/src/config"
	"github.com/Dataman-Cloud/go-template/src/db"
	"github.com/Dataman-Cloud/go-template/src/notification"
	//log "github.com/Sirupsen/logrus"
)

//var config *Config

var (
	envFile = flag.String("config", "env_file", "")
)

func main() {
	flag.Parse()

	config := config.InitConfig(*envFile)
	//log.SetLevel(config.LogLevel)
	db.MysqlInit()
	fmt.Println(config)
	message := notification.NewMessage()

	message.Type = "APP_CREATION"
	message.ResourceId = 44444444
	message.ResourceType = "3333333"
	//message.DumpBegin = time.Now().Add(time.Minute * -18)

	//	message.Remove()
	//	message.Persist()
	/*使用log示例*/
	//	LogDemo()

	//msgs := notification.LoadMessagesBefore(time.Minute * 50)

	//fmt.Println(msgs)

	//msgs = notification.LoadMessagesAfter(time.Minute * 50)

	//	fmt.Println(msgs)
	//notification.CleanOutdateMessagesBefore(time.Minute * 50)
	//	msgs := notification.LoadMessagesBefore(time.Minute * 50)

	//	fmt.Println(msgs)

	//	msgs = notification.LoadMessagesAfter(time.Minute * 50)

	//	fmt.Println(msgs)
	//	notification.CleanTooOldMessage(time.Minute * 50)

	engin := notification.NewEngine()

	go engin.Start()

	engin.Write(message)

	for {
		select {
		case <-time.After(time.Minute):
			fmt.Println("one")
		}
	}

	fmt.Println("birth cry")
}
