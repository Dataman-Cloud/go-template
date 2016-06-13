# Template

1、环境变量
   在程序初始化时，执行initEnv，如果找不到需要的环境变量，程序退出
   如果不同模块有自己的mysql，需要在黄金变量前加上模块名称


2、log打印


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
