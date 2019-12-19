package rpc

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/weikaishio/databus_kafka/common/conf/env"
	"github.com/weikaishio/databus_kafka/common/naming"
	"github.com/weikaishio/databus_kafka/common/naming/discovery"
	xtime "github.com/weikaishio/databus_kafka/common/time"

	. "github.com/smartystreets/goconvey/convey"
)

var c = &discovery.Config{
	Nodes:  []string{"api.domain.co"},
	Zone:   "sh001",
	Env:    "test",
	Key:    "Key",
	Secret: "Secret",
	Host:   "host_1",
}

var in = &naming.Instance{
	AppID:   "test2",
	Version: "1",
	Metadata: map[string]string{
		"test":    "1",
		"weight":  "8",
		"color":   "",
		"cluster": "red",
	},
}

var in2 = &naming.Instance{
	AppID:   "test3",
	Version: "1",
	Metadata: map[string]string{
		"test":    "1",
		"weight":  "8",
		"color":   "",
		"cluster": "red",
	},
}

var (
	svrAddr1, svrAddr2, svrAddr3 string
	once1, once2, once3          sync.Once
)

type TestArgs struct {
	A, B int
}

type TestReply struct {
	C int
}

type TestTimeout struct {
	T time.Duration
}

type TestRPC int

func startTestServer1() {
	svr := newServer()
	svr.RegisterName("RPC", new(TestRPC))
	var l net.Listener
	l, svrAddr1 = listenTCP()
	go svr.Accept(l)
}

func TestDiscoveryCli(t *testing.T) {
	env.Hostname = "host_1"
	env.Zone = "sh001"
	once1.Do(startTestServer1)

	Convey("test discovery cli", t, func() {
		once1.Do(startTestServer1)
		in.Addrs = []string{scheme + "://" + svrAddr1}
		dis := discovery.New(c)
		_, err := dis.Register(context.TODO(), in)
		So(err, ShouldBeNil)
		cli := NewDiscoveryCli("test2", &ClientConfig{
			Cluster: "",
			Timeout: xtime.Duration(time.Second),
		})
		time.Sleep(time.Second * 2)
		args := &TestArgs{7, 8}
		reply := new(TestReply)
		err = cli.Call(context.TODO(), "RPC.Add", args, reply)
		So(err, ShouldBeNil)
	})
	Convey("test discovery no zone", t, func() {
		env.Zone = "test2"
		cli := NewDiscoveryCli("test2", &ClientConfig{
			Cluster: "",
			Timeout: xtime.Duration(time.Second),
		})
		time.Sleep(time.Second * 2)
		args := &TestArgs{7, 8}
		reply := new(TestReply)
		err := cli.Call(context.TODO(), "RPC.Add", args, reply)
		So(err, ShouldBeNil)
	})

	Convey("test discovery with color", t, func() {
		env.Zone = "test2"
		cli := NewDiscoveryCli("test2", &ClientConfig{
			Color:   "red",
			Timeout: xtime.Duration(time.Second),
		})
		time.Sleep(time.Second * 2)
		args := &TestArgs{7, 8}
		reply := new(TestReply)
		err := cli.Call(context.TODO(), "RPC.Add", args, reply)
		So(err, ShouldBeNil)
	})
	Convey("test discovery with cluster", t, func() {
		env.Zone = "test2"
		cli := NewDiscoveryCli("test2", &ClientConfig{
			Cluster: "red",
			Timeout: xtime.Duration(time.Second),
		})
		time.Sleep(time.Second * 2)
		args := &TestArgs{7, 8}
		reply := new(TestReply)
		err := cli.Call(context.TODO(), "RPC.Add", args, reply)
		So(err, ShouldBeNil)
	})

	Convey("test conf Zone cli", t, func() {
		env.Zone = "testsh"
		once1.Do(startTestServer1)
		in2.Addrs = []string{scheme + "://" + svrAddr1}
		dis := discovery.New(c)
		_, err := dis.Register(context.TODO(), in2)
		So(err, ShouldBeNil)
		env.Zone = "sh001"
		cli := NewDiscoveryCli("test3", &ClientConfig{
			Cluster: "",
			Timeout: xtime.Duration(time.Second),
			Zone:    "testsh",
		})
		time.Sleep(time.Second * 2)
		args := &TestArgs{7, 8}
		reply := new(TestReply)
		err = cli.Call(context.TODO(), "RPC.Add", args, reply)
		So(err, ShouldBeNil)
	})
}
