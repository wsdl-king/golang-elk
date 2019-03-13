# golang-elk
go语言实现的elk(单节点)


### 依赖的软件

| 软件 | 版本|  
|:---------|:-------:|
| elasticSearch | 5.6.11 |
| kibana  |  5.6.11 |
| golang  |  1.10.4  |
| kafka  |  2.11-1.1.1 |
| zookeeper  |  3.4.9 |
| etcd  |  3.3.12 |
| docker  |  18.06.1-ce |

###项目分布
golang-elk ----目前只有一个main方法,etcd客户端,我用于实现分布式kv的修改<br>
logagent  ----日志收集客户端,实现功能为:收集现有各个服务器日志,发送到kafka<br>
logtransfer ---- 日志收集服务端,实现功能为:读取kafka日志,发送到elasticSearch<br>


<h3>在你所在的GOPATH进行go get操作,提取必要的包<h3>
### 依赖的package
######日志
go get github.com/astaxie/beego/logs<br>
######conig加载配置文件
go get github.com/pythonsite/config<br>
######etcd-client(v3)
go get go.etcd.io/etcd/clientv3<br>
######kafka-client
go get github.com/Shopify/sarama<br>
######实时读取文件的tail
github.com/hpcloud/tail<br>
######elasticSearch客户端 ps:前几天开源了官方的go-elasticSearch,有兴趣的小伙伴可以试一下.
gopkg.in/olivere/elastic.v2<br>

###使用docker进行服务安装
我的用户目录是/homg/qiwenshuai<br>
我的IP是192.168.88.152<br>
----
安装并运行elasticSearch 5.6.11<br>
docker pull elasticsearch:5.6.11<br>
mkdir -p /home/qiwenshuai/elasticsearch/config<br>
mkdir -p /home/qiwenshuai/elasticsearch/data<br>
echo "http.host: 0.0.0.0" >> /home/qiwenshuai/elasticsearch/config/elasticsearch.yml<br>
echo "http.cors.enabled: true" >> /home/qiwenshuai/elasticsearch/config/elasticsearch.yml<br>
echo "http.cors.allow-origin: "*"" >> /home/qiwenshuai/elasticsearch/config/elasticsearch.yml<br>
docker run --name elasticsearch5.6.11 -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -v /home/qiwenshuai/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml -v /home/qiwenshuai/elasticsearch/data:/usr/share/elasticsearch/data -d  elasticsearch:5.6.11<br>

安装并运行kibana 5.6.11<br>
docker pull kibana:5.6.11<br>
docker run --name kibana  -e ELASTICSEARCH_URL=elasticsearch5.6.11:9200 -p 5601:5601 -d kibana:5.6.11<br>

安装并运行zookeeper<br>
docker pull wurstmeister/zookeeper<br>
docker run  --name zookeeper   -p 2181:2181  -v /home/qiwenshuai/kafka/zoolog:/opt/zookeeper/data -d zookeeper<br>

安装并运行kafka<br>
docker pull wurstmeister/kafka
docker run   --name kafka -p  9092:9092 \
--link zookeeper \
--env KAFKA_ZOOKEEPER_CONNECT=192.168.88.152:2181 \
--env KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://192.168.88.152:9092 \
--env KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
 -d kafka
 
 安装etcd,这里我没有用docker安装,我一直pull不下来镜像,可能被墙了<br>
 ETCD_VER=v3.3.12<br>
 GOOGLE_URL=https://github.com/etcd-io/etcd/releases/download<br>
 DOWNLOAD_URL=${GOOGLE_URL}<br>
 rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz<br>
 rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test<br>
 curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz<br>
 tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1<br>
 rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz<br>
 查看版本<br>
 /tmp/etcd-download-test/etcd --version<br>
 ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl version<br>
 运行命令<br>
 nohup ./etcd --advertise-client-urls 'http://0.0.0.0:2379' --listen-client-urls 'http://0.0.0.0:2379' >log.out &<br>
 测试<br>
 ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl --endpoints=192.168.88.152:2379 put foo bar<br>
 ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl --endpoints=192.168.88.152:2379 get foo<br>
 

###配置修改
1.logtransfer
log_path=服务日志路径<br>
log_level=日志级别<br>
kafka_addr=kafka地址<br>
kafka_thread_num=kafka协程数量<br>
etcd_addr=etcd地址<br>
etcd_timeout = etcd超时时间,v3默认也是5s<br>
etcd_transfer_key = /logtransfer/%s/log_config (etcd-client设置的key,%s在我自己的服务中会用ip替换)<br>
es_addr=es地址<br>
es_thread_num= es协程数量<br>
2.logagent
log_path=服务日志路径<br>
log_level=日志级别<br>
kafka_addr=kafka地址<br>
kafka_thread_num=kafka协程数量<br>
etcd_addr=etcd地址<br>
etcd_timeout=etcd超时时间,v3默认也是5s<br>
etcd_watch_key=/logagent/%s/log_config(etcd-client设置的key,%s在我自己的服务中会用ip替换)<br>
3.golang-elk<br>
修改logconf(对应logagent的读取修改)<br>
// 这里介绍下下面的四个属性<br> 
// logconf 用于发送给etcd,然后logagent从chan里获取配置信息.<br>
// 1.topic: 这里用作于向kafka发送的topic<br>
// 2.log_path: 这里是我读取的日志路径<br>
// 3.service: 标示一个服务名称<br>
// 4.send_rate:发送速率,类似于tps的概念<br>
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
<br>
//修改transconf(对应logtransfer读取对应的topic)<br>
var transconf = `
[
    "eslservice_log"
]
`
<br>
###运行
到/go/src/logagent和/go/src/logtransfer下<br>
go build .<br>
./logagent ./logtransfer<br>
到 /go/src/golang-elk下<br>
go build .<br>
./main<br>
登录localhost:5601就可以看到kibana里的数据<br>
 
项目概述:<br>
golang语言实现的ELK,目前我司使用docker部署项目,微服务较多,查看日志文件就需要去服务器里tail或者docker logs ..等等噼里啪啦
然后我最近也在看golang,比较有兴趣,参照网上大神的ELK-fork的项目github地址:(https://github.com/pythonsite),进行必要的修改和注释
和部分中间件的安装.<br>
项目详细设计:<br>
1.golang-elk:<br>
    --只有一个main方法,后期我会继续提交,使用gogin,gorm做成restful接口.<br>
2.logagent:<br>
    --config.go :读取配置文件<br>
    --data.go: 结构体实例,etcd存储的value<br>
    --etcd.go: 初始化etcd,获取log配置,检测k的变化,向logConfChan发送数据<br>
    --ip.go: 用来得到本机内/外网ip<br>
    --kafka.go: 初始化kafka-producer,发送数据<br>
    --limit.go: 用于限制tps,每秒的限速<br>
    --main.go: 初始化log和config并且启动程序<br>
    --server.go: 从chan里加载配置,删除过期配置,使用tail读取文件发送到kafka<br>
3.logtransfer:<br>
    --es.go: 从msgchan读取数据,发送到es,里面有一个topicConfChan,etcd中的key发生变化以后进行修改消费者<br>
    --etcd.go: 初始化etcd,获取log配置,检测k的变化,向topicConfChan发送数据<br>
    --ip.go: 用来得到本机内/外网ip<br>
    --kafka.go: 初始化kafka-cousumer,消费数据发往msgchan<br>
    --main.go: 初始化log和config并且启动程序<br>
    
    
####扩展
    我决定将logagent构建成image放入各个想要收集服务日志的服务器中,然后进行pull
    logagent一定是无状态的,运行的时候只需要-v 指定好日志的映射目录,但是这样有一个问题,
    我通过golang-elk的client设置完log日志后,logagent读取的key是根据服务的ip来确定唯一性的
    所以我需要使用docker-swarm 或者k8s这样的容器管理工具确定指定的ip才可以.(目前未实现.笑哭.jpg)
    

鸣谢:pythonsite<br>
    -github 传送门: https://github.com/pythonsite<br>
    -csdn   传送门: http://www.cnblogs.com/zhaof/<br>