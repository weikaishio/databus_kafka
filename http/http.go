package http

import (
	"github.com/weikaishio/databus_kafka/common/log_b"
	bm "github.com/weikaishio/databus_kafka/common/net/http/blademaster"
	"github.com/weikaishio/databus_kafka/example_svr/conf"
	"github.com/weikaishio/databus_kafka/auth_service"
	"github.com/weikaishio/databus_kafka/tcp"
)

var (
	svc      *auth_service.Service
)

// Init init https
func Init(c *conf.Config, s *auth_service.Service) {
	svc = s
	// router
	router := bm.DefaultServer(c.HTTPServer)
	initRouter(router)
	// init internal server
	if err := router.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initRouter init local router api path.
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	e.GET("/databus/consumer/addrs", consumerAddrs)
	e.POST("/databus/pub", pub)
}

// ping check server ok
func ping(c *bm.Context) {
}

// register provid for discovery.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{
		"data": struct{}{},
	}, nil)
}

// consumerAddrs get consumer addrs.
func consumerAddrs(c *bm.Context) {
	group := c.Request.Form.Get("group")
	c.JSON(tcp.ConsumerAddrs(group))
}
