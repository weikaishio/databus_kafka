package direct

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/weikaishio/databus_kafka/common/net/netutil/breaker"
	"github.com/weikaishio/databus_kafka/common/net/rpc/warden"
	pb "github.com/weikaishio/databus_kafka/common/net/rpc/warden/proto/testproto"
	"github.com/weikaishio/databus_kafka/common/net/rpc/warden/resolver"
	xtime "github.com/weikaishio/databus_kafka/common/time"
)

type testServer struct {
	name string
}

func (ts *testServer) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: ts.name, Success: true}, nil
}

func (ts *testServer) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	panic("not implement error")
}

func createServer(name, listen string) *warden.Server {
	s := warden.NewServer(&warden.ServerConfig{Timeout: xtime.Duration(time.Second)})
	ts := &testServer{name}
	pb.RegisterGreeterServer(s.Server(), ts)
	go func() {
		if err := s.Run(listen); err != nil {
			panic(fmt.Sprintf("run warden server fail! err: %s", err))
		}
	}()
	return s
}

func TestMain(m *testing.M) {
	resolver.Register(New())
	ctx := context.TODO()
	s1 := createServer("server1", "127.0.0.1:18081")
	s2 := createServer("server2", "127.0.0.1:18082")
	defer s1.Shutdown(ctx)
	defer s2.Shutdown(ctx)
	os.Exit(m.Run())
}

func createTestClient(t *testing.T, connStr string) pb.GreeterClient {
	client := warden.NewClient(&warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	})
	conn, err := client.Dial(context.TODO(), connStr)
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	return pb.NewGreeterClient(conn)
}

func TestDirect(t *testing.T) {
	cli := createTestClient(t, "direct://default/127.0.0.1:18083,127.0.0.1:18082")
	count := 0
	for i := 0; i < 10; i++ {
		if resp, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
			t.Fatalf("TestDirect: SayHello failed!err:=%v", err)
		} else {
			if resp.Message == "server2" {
				count++
			}
		}
	}
	if count != 10 {
		t.Fatalf("TestDirect: get server2 times must be 10")
	}
}
