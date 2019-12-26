package redis_store

import (
	"github.com/weikaishio/databus_kafka/auth_service/model"
	"github.com/weikaishio/redis_orm"
)

type AuthRedis struct {
	redisOrm *redis_orm.Engine
}

func NewAuthRedis(redisOrm *redis_orm.Engine) *AuthRedis {
	return &AuthRedis{redisOrm: redisOrm}
}

func (d *AuthRedis) Auth(group string) (tb *model.Auth, has bool, err error) {
	authTb := &model.AuthTb{}
	has, err = d.redisOrm.GetByCondition(authTb, redis_orm.NewSearchConditionV2(group, group, "Group"))
	if err != nil || !has {
		return
	}
	appTb := &model.AppTb{Id: int64(authTb.AppId)}
	has, err = d.redisOrm.Get(appTb)
	if err != nil || !has {
		return
	}
	tb = &model.Auth{
		Group:     group,
		Topic:     authTb.Topic,
		Operation: authTb.Operation,
		Key:       appTb.AppKey,
		Secret:    appTb.AppSecret,
		Batch:     0,
		Cluster:   appTb.Cluster,
	}
	return
}
