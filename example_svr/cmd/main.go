package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weikaishio/databus_kafka/common/log_b"
	"github.com/weikaishio/databus_kafka/example_svr/conf"
	"github.com/weikaishio/databus_kafka/auth_service"
	"github.com/weikaishio/databus_kafka/http"
	"github.com/weikaishio/databus_kafka/tcp"
)
/*
CREATE TABLE `auth` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_id` INT(11) NOT NULL,
  `group_name` VARCHAR(45) NULL,
  `operation`TINYINT(4) DEFAULT 0,
  `topic` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `app` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_key` VARCHAR(45) NOT NULL,
  `app_secret` VARCHAR(64) NOT  NULL,
  `cluster` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`));


CREATE TABLE `topic` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `topic` VARCHAR(45) NULL,
  `cluster` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));

1、监听端口，auth，拿到cluster,topic等
2、接收cmd，执行pub or sub 然后写入
取队列数据校验权限，然后通过长链来取数据

用redis协议收发命令，x，好方式啊。。
redis-cli -h 127.0.0.1 -p 6205
pub: auth key:secret@group/topic=test156&role=pub
set 1 "xx"
sub: auth key:secret@group/topic=test156&role=sub
mget 1

log?
*/
func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("databus start")
	// service init
	svc := auth_service.New(conf.Conf.MySQL)
	http.Init(conf.Conf, svc)
	tcp.Init(conf.Conf.Addr, conf.Conf.Clusters, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("databus get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("databus exit")
			tcp.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
