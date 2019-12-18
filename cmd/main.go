package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weikaishio/databus_kafka/conf"
	"github.com/weikaishio/databus_kafka/service"
	"github.com/weikaishio/databus_kafka/tcp"
	"github.com/weikaishio/databus_kafka/common/log_b"
)
/*
CREATE TABLE `auth` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_id` INT(11) NULL,
  `group_name` VARCHAR(45) NULL,
  `operation` VARCHAR(45) NULL,
  `topic` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `app` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_key` VARCHAR(45) NULL,
  `app_secret` VARCHAR(64) NULL,
  `cluster` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `auth2` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_id` INT(11) NULL,
  `group` VARCHAR(45) NULL,
  `operation` TINYINT(4) NULL,
  `app_key` VARCHAR(45) NULL,
  `app_secret` VARCHAR(64) NULL,
  `number` VARCHAR(45) NULL,
  `is_delete` TINYINT(4) NULL,
  `topic_id` INT(11) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `topic` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `topic` VARCHAR(45) NULL,
  `cluster` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));
CREATE TABLE `app2` (44
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_secret` VARCHAR(64) NULL,
  `app_key` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));

1、监听端口，auth，拿到cluster,topic等
2、接收cmd，执行pub or sub 然后写入
没意义。。。 取队列数据校验权限，然后通过长链来取数据
才用redis协议收发命令，x，好方式啊。。
redis-cli -h 127.0.0.1 -p 6205

比如app/admin/main/member/dao/member.go
func (d *Dao) PubExpMsg(ctx context.Context, msg *model.AddExpMsg) error {
	return d.expMsgDatabus.Send(ctx, strconv.FormatInt(msg.Mid, 10), msg)
}
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
	svc := service.New(conf.Conf)
	//http.Init(conf.Conf, svc)
	tcp.Init(conf.Conf, svc)
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
