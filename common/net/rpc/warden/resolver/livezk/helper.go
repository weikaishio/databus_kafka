package livezk

import (
	"context"
	"fmt"
	"net"

	"github.com/weikaishio/databus_kafka/common/naming"
	lz "github.com/weikaishio/databus_kafka/common/naming/livezk"
	"github.com/weikaishio/databus_kafka/common/net/ip"
)

// Register self grpc service to live zookeeper
func Register(config *lz.Zookeeper, addr string, discoveryID string) (context.CancelFunc, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	z, err := lz.New(config)
	if err != nil {
		return nil, err
	}
	internalIP := ip.InternalIP()
	ins := &naming.Instance{
		AppID: discoveryID,
		Addrs: []string{fmt.Sprintf("grpc://%s:%s", internalIP, port)},
	}
	return z.Register(context.Background(), ins)
}
