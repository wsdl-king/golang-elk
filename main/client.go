package main

//一个etcd-client 我用作于分布式配置管理,后期我会直接跟go-gin+gorm结合,使用rest进行配置
import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
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
        "log_path":"/home/qiwenshuai/logs/aaa.log,/home/qiwenshuai/logs/da.log",
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
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.88.152:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed,err:", err)
		return
	}
	fmt.Println("connect success")
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/logagent/192.168.88.152/log_config", logconf)
	_, err = cli.Put(ctx, "/logtransfer/192.168.88.152/log_config", transconf)
	cancel()
	if err != nil {
		fmt.Println("put failed ,err:", err)
		return
	}
}
