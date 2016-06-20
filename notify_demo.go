package main

import (
	"github.com/Dataman-Cloud/go-template/src/notification"
)

/* nofitication组件功能：在发生事件时，可以把消息发送给注册的组件地址，可以通过配置mode，来指定消息失败是否重传，保证订阅方收到消息。
1、通过配置 Notification环境变量
       config.Notification = `http://192.168.1.106:8080/postMan?name=li&mode=strict&notification_types=APP_CREATION|
						http://127.0.0.1:8089/postMan?name=fang&mode=strict&notification_types=APP_CREATION`
   通过url格式进行配置，不同的sink之间通过 | 分割，
   url:   http://192.168.1.106:8080/postMan    消息要发送到的地址
   name： sink的名字 一般是订阅方的名字
   mode： 模式有两种 strict：该模式下消息发送失败会进行重传，尝试多次失败后会存到数据库，之后定期发送，保证消息送达
   		     best_deliver：该模式消息发送失败不会重新发送
    notification_types: 关心的事件  目前只有两种 APP_DELETION 应用删除    APP_CREATION 应用创建

2、strict模式下几个关键时间点
	1)、消息发送失败后，会根据RetryAfer来进行重传，再失败后会再RetryAfter*RetryBackoffFactor后进行重传，以此类推。。。。
		当时间超过MessageStaleTime后，消息暂停重传。
	2)、会按MessageGCTime周期来进行数据库里消息的清理和重传    超过MessageDeleteTime时间的将被删除，其他消息会重传一次，失败后等下一个周期来触发
	3)、程序启动时会先在数据库中查找 MessageStaleTime时间内的消息，进行发送

	目前在notifycation.go 里定义，使用时，可以通过环境变量传入
		MessageStaleTime   = time.Duration(time.Minute * 10) // message not sent after this value marked as stale
		RetryAfer          = time.Duration(time.Second * 2)  // initial delay if message not sent successfully
		RetryBackoffFactor = 2                               // 每次重传后的backoff值
		MessageGCTime      = time.Minute * 10                // 每隔一段时间清理下数据库里发送失败的消息
		MessageDeleteTime  = time.Hour * 24                  // 数据库里超过该时间的数据将被彻底删除


3、数据库表创建在sql目录下   0001_notification_create_database.up.sql
*/

func NotifyDemo() {

	// 启动nofification
	engin := notification.NewEngine()

	go engin.Start()

	//启动后就可以通过 Write方法进行消息发送

	message := notification.NewMessage()

	message.Type = "APP_CREATION" //指定消息类型
	message.ResourceId = 44444444
	message.ResourceType = "APP" //目前只是APP  为了可能还有其他类型 方便以后扩展

	engin.Write(message)

	//可以通过Stop方法停止engine
	engin.Stop()

}
