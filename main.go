package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

/* 如果不同模块需要配置各自的mysql 前面加上模块名 方便运维配置*/
type EnvConfig struct {
	Env1           string
	Env2           string
	A_MYSQL_PASSWD string
	LOG_LEVEL      string
}

var Env EnvConfig

/* 缺少环境变量 程序退出 使用log.Fatalf打印*/
func exitMissingEnv(env string) {
	log.Fatalf("exit missing env %s", env)
}

func initEnv() {
	Env1 := os.Getenv("ENV1")
	if Env1 == "" {
		exitMissingEnv("ENV1")
	}

	Env2 := os.Getenv("ENV2")
	if Env2 == "" {
		exitMissingEnv("ENV2")
	}

	A_MYSQL_PASSWD := os.Getenv("A_MYSQL_PASSWD")
	if A_MYSQL_PASSWD == "" {
		exitMissingEnv("A_MYSQL_PASSWD")
	}

	LOG_LEVEL := os.Getenv("LOG_LEVEL")

	Env.Env1 = Env1
	Env.Env2 = Env2
	Env.A_MYSQL_PASSWD = A_MYSQL_PASSWD
	Env.LOG_LEVEL = LOG_LEVEL
}

func main() {
	initEnv()
	if Env.LOG_LEVEL == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	}
	/*使用log示例*/
	Test()
	fmt.Println("birth cry")

}
