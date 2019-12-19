package conf

import (
	"flag"
	"fmt"

	"github.com/weikaishio/databus_kafka/common/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	//client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	//BM            *bm.ServerConfig
	//Tracer        *trace.Config

	//HTTPClient    *HTTPClient

	ExampleSubMsgDatabus *databus.Config
	ExamplePubMsgDatabus *databus.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "./config.toml", "config path")
}

//Init int config
func Init() error {
	fmt.Printf("confPath: %s\n",confPath)
	if confPath != "" {
		return local()
	}
	return fmt.Errorf("confPath is nil")
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}