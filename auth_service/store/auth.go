package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/weikaishio/databus_kafka/auth_service/model"
	"github.com/weikaishio/databus_kafka/auth_service/store/db_store"
	"github.com/weikaishio/databus_kafka/auth_service/store/redis_store"
	"github.com/weikaishio/databus_kafka/common/database/sql"
	tim "github.com/weikaishio/databus_kafka/common/time"
	"github.com/weikaishio/redis_orm"
)

type AuthStore struct {
	Dao    *db_store.Dao
	Auth   func(group string) (tb *model.Auth, has bool, err error)
	AuthDB func(c context.Context) (auths map[string]*model.Auth, err error)
}

type Options struct {
	Addr               string
	Password           string
	DB                 int
	DialTimeout        tim.Duration
	ReadTimeout        tim.Duration
	WriteTimeout       tim.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         tim.Duration
	PoolTimeout        tim.Duration
	IdleTimeout        tim.Duration
	IdleCheckFrequency tim.Duration
}

func LoadAuthDBStore(mysql *sql.Config) *AuthStore {
	dao := db_store.New(mysql)
	authStore := &AuthStore{
		AuthDB: dao.Auth,
	}
	authStore.Dao = dao
	return authStore
}
func LoadAuthRedisStore(redisOpt *Options) (*AuthStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:               redisOpt.Addr,
		Password:           redisOpt.Password,
		DB:                 redisOpt.DB,
		DialTimeout:        time.Duration(redisOpt.DialTimeout),
		ReadTimeout:        time.Duration(redisOpt.ReadTimeout),
		WriteTimeout:       time.Duration(redisOpt.WriteTimeout),
		IdleTimeout:        time.Duration(redisOpt.IdleTimeout),
		IdleCheckFrequency: time.Duration(redisOpt.IdleCheckFrequency),
		MaxConnAge:         time.Duration(redisOpt.MaxConnAge),
		MinIdleConns:       redisOpt.MinIdleConns,
		PoolSize:           redisOpt.PoolSize,
		PoolTimeout:        time.Duration(redisOpt.PoolTimeout),
	})
	ping, err := client.Ping().Result()
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	if strings.ToLower(ping) != "pong" {
		return nil, errors.New("redis failed to ping")
	}
	redisOrm := redis_orm.NewEngine(client)
	authRedis := redis_store.NewAuthRedis(redisOrm)
	return &AuthStore{
		Auth: authRedis.Auth,
	}, nil
}
