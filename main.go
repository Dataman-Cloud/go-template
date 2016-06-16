package main

import (
	"flag"
	"fmt"

	"github.com/Dataman-Cloud/go-template/src/config"
	"github.com/Dataman-Cloud/go-template/src/db"
	"github.com/Dataman-Cloud/go-template/src/notification"
	log "github.com/Sirupsen/logrus"
)

//var config *Config

var (
	envFile = flag.String("config", "env_file", "")
)

func main() {
	flag.Parse()

	config := config.InitConfig(*envFile)
	log.SetLevel(config.LogLevel)
	db.MysqlInit()
	fmt.Println(config)
	message := notification.NewMessage()

	message.Id = "11111111"
	message.Type = "strict"
	message.ResourceId = "22222222"
	message.ResourceType = "3333333"

	message.Persist()
	/*使用log示例*/
	//	LogDemo()

	fmt.Println("birth cry")
}
