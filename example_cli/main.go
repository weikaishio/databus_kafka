package main

import (
	"context"
	"fmt"

	log "github.com/weikaishio/databus_kafka/common/log_b"

	"github.com/weikaishio/databus_kafka/example_cli/conf"

	"github.com/weikaishio/databus_kafka/common/queue/databus"
)

func main() {
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	fmt.Printf("conf:%v,%v\n", conf.Conf, conf.Conf.ExamplePubMsgDatabus)

	examplePubMsgDatabus := databus.New(conf.Conf.ExamplePubMsgDatabus)
	err := examplePubMsgDatabus.Send(context.Background(), "aaaa", "example publish")
	if err != nil {
		fmt.Printf("err:%v", err)
	}

	exampleSubMsgDatabus := databus.New(conf.Conf.ExampleSubMsgDatabus)
	for {
		select {
		case msg := <-exampleSubMsgDatabus.Messages():
			fmt.Printf("msg:%v,%v,%v,%v,%v,%v\n", msg.Key, string(msg.Value), msg.Topic, msg.Offset, msg.Partition, msg.Timestamp)
		}
	}
}
