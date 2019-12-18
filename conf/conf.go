package conf

import (
	"errors"
	"flag"
	"github.com/weikaishio/databus_kafka/common/database/sql"

	bm "github.com/weikaishio/databus_kafka/common/net/http/blademaster"

	"github.com/weikaishio/databus_kafka/common/conf"
	"github.com/weikaishio/databus_kafka/common/log_b"

	"github.com/BurntSushi/toml"
)

var (
	// Conf global config variable
	Conf     = &Config{}
	confPath string
	client   *conf.Client
)

// Config databus config struct
type Config struct {
	// base
	Addr     string
	Clusters map[string]*Kafka
	// Log
	Log *log.Config
	// http
	HTTPServer *bm.ServerConfig
	// mysql
	MySQL *sql.Config
}

// Kafka contains cluster, brokers, sync.
type Kafka struct {
	Cluster string
	Brokers []string
}

func init() {
	flag.StringVar(&confPath, "conf", "/Users/wangjiangang/beta/workspace/golang/src/github.com/weikaishio/databus_kafka/common/cmd/databus-test.toml", "config path")
}

//Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}