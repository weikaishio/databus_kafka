package db_store

import (
	"context"

	"github.com/weikaishio/databus_kafka/auth_service/model"

	"github.com/weikaishio/databus_kafka/common/log_b"
)

const (
	_getAuthSQL = `SELECT auth.group_name,auth.operation,app.app_key,app.app_secret,auth.topic,app.cluster
				FROM auth LEFT JOIN app ON auth.app_id=app.id`
)

// Auth verify group,topic,key
func (d *Dao) Auth(c context.Context) (auths map[string]*model.Auth, err error) {
	auths = make(map[string]*model.Auth)
	// auth
	rows, err := d.db.Query(c, _getAuthSQL)
	if err != nil {
		log.Error("getAuthStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Auth{}
		if err = rows.Scan(&a.Group, &a.Operation, &a.Key, &a.Secret, &a.Topic, &a.Cluster); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		auths[a.Group] = a
	}
	return
}
