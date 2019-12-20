package service

import (
	"github.com/weikaishio/databus_kafka/common/rate/limit/bench/stress/conf"
)

// Service struct
type Service struct {
	c *conf.Config
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
	}
	return s
}
