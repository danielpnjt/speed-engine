package bank

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"github.com/danielpnjt/speed-engine/internal/domain/repositories"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
	"github.com/danielpnjt/speed-engine/internal/pkg/types"
	"gorm.io/gorm"
)

type service struct {
	db             *gorm.DB
	bankRepository repositories.Bank
	redisWrapper   redis.Wrapper
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *gorm.DB) *service {
	s.db = db
	return s
}

func (s *service) SetBankRepository(repository repositories.Bank) *service {
	s.bankRepository = repository
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
	if s.bankRepository == nil {
		panic("bankRepository is nil")
	}
	if s.redisWrapper == nil {
		panic("redisWrapper is nil")
	}
	return s
}

func (s *service) SubmitBank(ctx context.Context, req SubmitBank) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	fmt.Println("===========", userData, userData.Username)
	userID := userData.ID

	bank := &entities.Bank{
		UserID:        userID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		BankName:      req.BankName,
	}
	err = s.bankRepository.Create(ctx, bank)
	if err != nil {
		slog.ErrorContext(ctx, "failed to submit bank", err)
		err = fmt.Errorf("failed to submit bank")
		return
	}
	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    "",
		Errors:  make([]string, 0),
	}
	return
}

func (s *service) FindAll(ctx context.Context) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	userID := userData.ID

	banks, err := s.bankRepository.FindByUserID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find bank by userID", err)
		err = fmt.Errorf("failed to find bank by userID")
		return
	}
	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    banks,
		Errors:  make([]string, 0),
	}
	return
}
