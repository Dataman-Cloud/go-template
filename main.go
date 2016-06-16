package main

import (
	"flag"
	"fmt"

	"github.com/Dataman-Cloud/go-template/src/config"
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

	fmt.Println(config)

	/*使用log示例*/
	LogDemo()

	fmt.Println("birth cry")
}
