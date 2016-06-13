package main

import (
	"bufio"
	"errors"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

/* 如果不同模块需要配置各自的mysql 前面加上模块名 方便运维配置*/
type Config struct {
	Env1           string
	Env2           string
	A_MYSQL_PASSWD string
	LogLevel       log.Level
	Arr            []string
}

func InitConfig(envFile string) *Config {
	loadEnvFile(envFile)

	config := &Config{}

	env1 := os.Getenv("ENV1")
	if env1 == "" {
		exitMissingEnv("ENV1")
	}

	env2 := os.Getenv("ENV2")
	if env2 == "" {
		exitMissingEnv("ENV2")
	}

	a_mysql_passwd := os.Getenv("A_MYSQL_PASSWD")
	if a_mysql_passwd == "" {
		exitMissingEnv("A_MYSQL_PASSWD")
	}

	log_level := os.Getenv("LOG_LEVEL")
	if log_level == "error" {
		config.LogLevel = log.ErrorLevel
	} else if log_level == "info" {
		config.LogLevel = log.InfoLevel
	} else {
		config.LogLevel = log.DebugLevel
	}

	arr := os.Getenv("ARR")
	if arr == "" {
		exitMissingEnv("ARR")
	}

	config.Env1 = env1
	config.Env2 = env2
	config.A_MYSQL_PASSWD = a_mysql_passwd
	config.Arr = strings.Split(arr, ",")

	return config
}

func loadEnvFile(envfile string) {
	// load the environment file
	f, err := os.Open(envfile)
	if err == nil {
		defer f.Close()

		r := bufio.NewReader(f)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}

			key, val, err := parseln(string(line))
			if err != nil {
				continue
			}

			if len(os.Getenv(strings.ToUpper(key))) == 0 {
				os.Setenv(strings.ToUpper(key), val)
			}
		}
	}
}

// helper function to parse a "key=value" environment variable string.
func parseln(line string) (key string, val string, err error) {
	line = removeComments(line)
	if len(line) == 0 {
		return
	}
	splits := strings.SplitN(line, "=", 2)

	if len(splits) < 2 {
		err = errors.New("missing delimiter '='")
		return
	}

	key = strings.Trim(splits[0], " ")
	val = strings.Trim(splits[1], ` "'`)
	return

}

// helper function to trim comments and whitespace from a string.
func removeComments(s string) (_ string) {
	if len(s) == 0 || string(s[0]) == "#" {
		return
	} else {
		index := strings.Index(s, " #")
		if index > -1 {
			s = strings.TrimSpace(s[0:index])
		}
	}
	return s
}

func exitMissingEnv(env string) {
	log.Fatalf("program exit missing config for env %s", env)
	os.Exit(1)
}
