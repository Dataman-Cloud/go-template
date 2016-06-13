package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

/* 如果不同模块需要配置各自的mysql 前面加上模块名 方便运维配置*/
type Config struct {
	Env1           string
	Env2           string
	A_MYSQL_PASSWD string
	LogLevel       log.Level
}

/* 缺少环境变量 程序退出 使用log.Fatalf打印*/
func exitMissingEnv(env string) {
	log.Fatal("exit missing value for env %s", env)
	os.Exit(1)
}

func InitConfig(envFile string) *Config {
	loadEnvFile(envFile)

	config := &Config{}

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

	config.Env1 = Env1
	config.Env2 = Env2
	config.A_MYSQL_PASSWD = A_MYSQL_PASSWD

	return config
}

func loadEnvFile(envfile string) {
}
