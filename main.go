package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var config *Config

var (
	envFile = flag.String("config", "env_file", "")
)

func main() {
	flag.Parse()

	config := InitConfig(*envFile)
	log.SetLevel(config.LogLevel)

	/*使用log示例*/
	LogDemo()

	fmt.Println("birth cry")
}
