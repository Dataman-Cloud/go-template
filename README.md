# GO项目配置和LOG约定

## 配置

1、从配置文件写入环境变量
  在程序初始化时，优先读取配置文件中配置，如环境变量中没有相应配置，
则写入环境变量.

2、从环境变量构造Config结构体


## Log

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
