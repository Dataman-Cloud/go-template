package main

import (
	log "github.com/Sirupsen/logrus"
)

/*

	log使用示例
	1、统一采用 github.com/Sirupsen/logrus
	2、该log有六种级别
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	// Calls os.Exit(1) after logging
	log.Fatal("Bye.")
	// Calls panic() after logging
	log.Panic("I'm bailing.")

	Fatal和Panic会退出程序 不建议使用

	3、发生错误需要用log.Error打印错误信息
	   正常流程使用log.Infoln打印，不建议打印太多
	      如果加用于调试的信息 加DEBUG log.Debugln打印

*/

func getInfo() string {

	return "iiiiii"
}

func doSomething() error {
	i := getInfo()

	/* 如果需要对信息进行记录 可以用log.Debug*/

	log.Debugln(i)

	return nil
}

func Test() {

	err := doSomething()

	/*发生错误需要用log.Error打印错误信息*/
	if err != nil {
		log.Errorln(err.Error())
	}

	/*可以根据情况进行正常流程打印*/
	log.Infoln("do_something successfully")

}
