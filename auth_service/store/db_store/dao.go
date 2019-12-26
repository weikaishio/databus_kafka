package db_store

import (
	"context"

	"github.com/weikaishio/databus_kafka/common/database/sql"
)

// Dao mysql struct
type Dao struct {
	db *sql.DB
}

// New new a Dao and return
func New(mysql *sql.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(mysql),
	}
	return
}

// Ping ping mysql
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}

// Close release mysql connection
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
