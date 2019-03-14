package main

import (
	"github.com/gin-gonic/gin"
	"golang-elk/route"
	"net/http"
	"time"
)

// 这里介绍下下面的四个属性
// logconf 用于发送给etcd,然后logagent从chan里获取配置信息.
// 1.topic: 这里用作于向kafka发送的topic
// 2.log_path: 这里是我读取的日志路径
// 3.service: 标示一个服务名称
// 4.send_rate:发送速率,类似于tps的概念
var logconf = `
[
    {
        "topic":"eslservice_log",
        "log_path":"/home/qiwenshuai/logs/dsda.log",
        "service":"eslservice",
        "send_rate":50000
    }
]
`
var transconf = `
[
    "eslservice_log"
]
`

func main() {
	engine := gin.New()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
		//如果不加单位 则是Nanosecond 纳秒
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	engine.POST("/add", route.Add())
	srv.ListenAndServe()
}
