package proxy_test

import (
	"github.com/weikaishio/databus_kafka/common/net/http/blademaster"
	"github.com/weikaishio/databus_kafka/common/net/http/blademaster/middleware/proxy"
)

// This example create several reverse proxy to show how to use proxy middleware.
// We proxy three path to `api.domain.com` and return response without any changes.
func Example() {
	proxies := map[string]string{
		"/index":        "http://api.domain.com/html/index",
		"/ping":         "http://api.domain.com/api/ping",
		"/api/versions": "http://api.domain.com/api/web/versions",
	}

	engine := blademaster.Default()
	for path, ep := range proxies {
		engine.GET(path, proxy.NewAlways(ep))
	}
	engine.Run(":18080")
}

// This example create several reverse proxy to show how to use jd proxy middleware.
// The request will be proxied to destination only when request is from specified datacenter.
func ExampleNewZoneProxy() {
	proxies := map[string]string{
		"/index":        "http://api.domain.com/html/index",
		"/ping":         "http://api.domain.com/api/ping",
		"/api/versions": "http://api.domain.com/api/web/versions",
	}

	engine := blademaster.Default()
	// proxy to specified destination
	for path, ep := range proxies {
		engine.GET(path, proxy.NewZoneProxy("sh004", ep), func(ctx *blademaster.Context) {
			ctx.String(200, "Origin")
		})
	}
	// proxy with request path
	ug := engine.Group("/update", proxy.NewZoneProxy("sh004", "http://sh001-api.domain.com"))
	ug.POST("/name", func(ctx *blademaster.Context) {
		ctx.String(500, "Should not be accessed")
	})
	ug.POST("/sign", func(ctx *blademaster.Context) {
		ctx.String(500, "Should not be accessed")
	})
	engine.Run(":18080")
}
