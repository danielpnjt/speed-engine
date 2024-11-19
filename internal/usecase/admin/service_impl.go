package admin

import (
	"github.com/danielpnjt/speed-engine/internal/domain/repositories"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"gorm.io/gorm"
)

type service struct {
	db                    *gorm.DB
	transactionRepository repositories.Transaction
	redisWrapper          redis.Wrapper
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *gorm.DB) *service {
	s.db = db
	return s
}

func (s *service) SetTransactionRepository(repository repositories.Transaction) *service {
	s.transactionRepository = repository
	return s
}

func (s *service) SetRedisWrapper(wrapper redis.Wrapper) *service {
	s.redisWrapper = wrapper
	return s
}

func (s *service) Validate() Service {
	if s.db == nil {
		panic("db is nil")
	}
	if s.transactionRepository == nil {
		panic("transactionRepository is nil")
	}
	if s.redisWrapper == nil {
		panic("redisWrapper is nil")
	}
	return s
}
