package auth_service

import (
	"context"
	"errors"
	"time"

	"github.com/weikaishio/databus_kafka/auth_service/model"
	"github.com/weikaishio/databus_kafka/auth_service/store"
	"github.com/weikaishio/databus_kafka/common/database/sql"
	"github.com/weikaishio/databus_kafka/common/stat/prom"

	"github.com/weikaishio/databus_kafka/common/log_b"
)

const (
	_authUpdateInterval = 1 * time.Minute
)

// Service service instance
type Service struct {
	authStore *store.AuthStore
	// auth
	auths map[string]*model.Auth
	// the auth of cluster changed
	clusterChan chan model.Auth
	// stats prom
	StatProm  *prom.Prom
	CountProm *prom.Prom
	TimeProm  *prom.Prom
}

// New new and return service
func New(authType int, mysql *sql.Config, redisOpt *store.Options) (s *Service, err error) {
	s = &Service{
		// cluster
		clusterChan: make(chan model.Auth, 5),
		// stats prom
		StatProm: prom.New().WithState("go_databus_state", []string{"role", "group", "topic", "partition"}),
		// count prom: count consumer and producer partition speed
		CountProm: prom.New().WithState("go_databus_counter", []string{"operation", "group", "topic"}),
		TimeProm:  prom.New().WithTimer("go_databus_timer", []string{"group"}),
	}
	if authType == 0 {
		s.authStore = store.LoadAuthDBStore(mysql)
		_ = s.fillAuth()
		go s.proc()
	} else {
		authStore, err := store.LoadAuthRedisStore(redisOpt)
		if err != nil {
			return nil, err
		}
		s.authStore = authStore
	}
	return
}

// Ping check mysql connection
func (s *Service) Ping(c context.Context) error {
	if s.authStore.Dao != nil {
		return s.authStore.Dao.Ping(c)
	}
	return nil
}

// Close close mysql connection
func (s *Service) Close() {
	if s.authStore.Dao != nil {
		s.authStore.Dao.Close()
	}
}

func (s *Service) proc() {
	for {
		_ = s.fillAuth()
		time.Sleep(_authUpdateInterval)
	}
}

func (s *Service) fillAuth() (err error) {
	if s.authStore.AuthDB == nil {
		return errors.New("s.authStore.AuthDB is nil")
	}
	auths, err := s.authStore.AuthDB(context.Background())
	if err != nil {
		log.Error("service.fillAuth error(%v)", err)
		return
	}
	var changed []*model.Auth
	// check cluster change event
	for group, nw := range auths {
		old, ok := s.auths[group]
		if !ok {
			continue
		}
		if old.Cluster != nw.Cluster {
			changed = append(changed, old)
			log.Info("cluster changed group(%s) topic(%s) oldCluster(%s) newCluster(%s)", old.Group, old.Topic, old.Cluster, nw.Cluster)
		}
	}
	s.auths = auths
	for _, ch := range changed {
		s.clusterChan <- *ch
	}
	return
}

// AuthApp check auth from cache
func (s *Service) AuthApp(group string) (a *model.Auth, ok bool) {
	if len(s.auths) > 0 {
		a, ok = s.auths[group]
	} else {
		a, ok, _ = s.authStore.Auth(group)
	}
	return
}

// ClusterEvent return cluster change event
func (s *Service) ClusterEvent() (group <-chan model.Auth) {
	return s.clusterChan
}
