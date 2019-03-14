package client

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"golang-elk/conf"
	"time"
)

func Put(key string, attr string) {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.NewConfig("./config.yml").EtcdAddr},
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	_, e2 := cli.Put(ctx, key, attr)
	if e2 != nil {
		panic(e2)
	}
	cancelFunc()
}
