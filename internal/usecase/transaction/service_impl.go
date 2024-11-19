package transaction

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"github.com/danielpnjt/speed-engine/internal/domain/repositories"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/payment"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
	"github.com/danielpnjt/speed-engine/internal/pkg/types"
	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"gorm.io/gorm"
)

type service struct {
	db                    *gorm.DB
	transactionRepository repositories.Transaction
	userRepository        repositories.User
	bankRepository        repositories.Bank
	redisWrapper          redis.Wrapper
	paymentWrapper        payment.Wrapper
	worker                *machinery.Server
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

func (s *service) SetUserRepository(repository repositories.User) *service {
	s.userRepository = repository
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

func (s *service) SetPaymentWrapper(wrapper payment.Wrapper) *service {
	s.paymentWrapper = wrapper
	return s
}

func (s *service) SetWorker(worker *machinery.Server) *service {
	s.worker = worker
	return s
}

func (s *service) Validate() Service {
	if s.db == nil {
		panic("db is nil")
	}
	if s.transactionRepository == nil {
		panic("transactionRepository is nil")
	}
	if s.userRepository == nil {
		panic("userRepository is nil")
	}
	if s.bankRepository == nil {
		panic("bankRepository is nil")
	}
	if s.redisWrapper == nil {
		panic("redisWrapper is nil")
	}
	if s.paymentWrapper == nil {
		panic("paymentWrapper is nil")
	}
	if s.worker == nil {
		panic("worker is nil")
	}
	return s
}

func (s *service) Generate(ctx context.Context, req GenerateRequest) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	reference, err := utils.GeneratePaymentRef(userData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate payment ref", err)
		err = fmt.Errorf("can't generate payment ref")
		return
	}

	vaRequest := payment.CreateVARequest{
		ExternalID:     reference,
		ExpectedAmount: req.Amount,
	}
	va, err := s.paymentWrapper.CreateVA(ctx, vaRequest)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate va", err)
		err = fmt.Errorf("can't generate va")
		return
	}

	entryData := &entities.Transaction{
		UserID:    userData.ID,
		BankID:    0,
		Amount:    va.Data.ExpectedAmount,
		Type:      "in",
		Reference: va.Data.ExternalID,
		Status:    va.Data.Status,
		ExpiredAt: *va.Data.ExpirationDate,
	}

	err = s.transactionRepository.Create(ctx, entryData)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate transfer", err)
		err = fmt.Errorf("can't generate transfer")
		return
	}

	eta := time.Now().UTC().Add(config.GetDuration("worker.speed_engine.delayFirstRetry"))
	signature := &tasks.Signature{
		Name: "enqueue-speed_engine-topup",
		Args: []tasks.Arg{
			{Name: "external", Type: "string", Value: va.Data.ExternalID},
		},
		RetryCount:   config.GetInt("worker.speed_engine.retryCount"),
		RetryTimeout: config.GetInt("worker.speed_engine.retryTimeout"),
		ETA:          &eta,
	}

	_, err = s.worker.SendTaskWithContext(ctx, signature)
	if err != nil {
		slog.ErrorContext(ctx, "send task worker", err)
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    va,
		Errors:  make([]string, 0),
	}
	return
}

func (s *service) TopUp(ctx context.Context, reference string) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	transaction, err := s.transactionRepository.FindByReference(ctx, reference)
	if err != nil {
		slog.ErrorContext(ctx, "can't get transaction", err)
		err = fmt.Errorf("can't get transaction")
		return
	}

	topUpRequest := payment.TopUpRequest{
		ExternalID: reference,
	}
	topUp, err := s.paymentWrapper.TopUp(ctx, topUpRequest)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate va", err)
		err = fmt.Errorf("can't generate va")
		return
	}

	if topUp.Data.Status != "COMPLETED" {
		eta := time.Now().UTC().Add(config.GetDuration("worker.speed_engine.delayFirstRetry"))
		signature := &tasks.Signature{
			Name: "enqueue-speed_engine-topup",
			Args: []tasks.Arg{
				{Name: "externalId", Type: "string", Value: reference},
			},
			RetryCount:   config.GetInt("worker.speed_engine.retryCount"),
			RetryTimeout: config.GetInt("worker.speed_engine.retryTimeout"),
			ETA:          &eta,
		}

		_, err = s.worker.SendTaskWithContext(ctx, signature)
		if err != nil {
			slog.ErrorContext(ctx, "send task worker", err)
		}
		return
	}

	transaction.Status = topUp.Data.Status
	err = s.transactionRepository.Create(ctx, &transaction)
	if err != nil {
		slog.ErrorContext(ctx, "can't create transfer", err)
		err = fmt.Errorf("can't create transfer")
		return
	}

	newUser, err := s.userRepository.FindByID(ctx, userData.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find user", err)
		err = fmt.Errorf("failed to find user")
		return
	}

	adminFee := float64(4500)
	newUser.Balance += transaction.Amount - adminFee
	err = s.userRepository.Create(ctx, &newUser)
	if err != nil {
		slog.ErrorContext(ctx, "failed to modify balance", err)
		err = fmt.Errorf("failed to modify balance")
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

func (s *service) Withdraw(ctx context.Context, req WithdrawRequest) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	reference, err := utils.GeneratePaymentRef(userData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate payment ref", err)
		err = fmt.Errorf("can't generate payment ref")
		return
	}

	bank, err := s.bankRepository.FindByID(ctx, req.BankID)
	if err != nil {
		slog.ErrorContext(ctx, "can't get bank", err)
		err = fmt.Errorf("can't get bank")
		return
	}

	withdrawRequest := payment.WithdrawRequest{
		ExternalID:        reference,
		BankCode:          bank.BankName,
		AccountHolderName: bank.AccountName,
		AccountNumber:     bank.AccountNumber,
		Amount:            req.Amount,
		Description:       "disbursement to user",
	}

	withdraw, err := s.paymentWrapper.Withdraw(ctx, withdrawRequest)
	if err != nil {
		slog.ErrorContext(ctx, "can't generate va", err)
		err = fmt.Errorf("can't generate va")
		return
	}

	transaction := entities.Transaction{
		UserID:    userData.ID,
		Amount:    float64(withdraw.Data.Amount),
		Type:      "out",
		Reference: reference,
		Status:    withdraw.Data.Status,
	}
	err = s.transactionRepository.Create(ctx, &transaction)
	if err != nil {
		slog.ErrorContext(ctx, "can't create transfer", err)
		err = fmt.Errorf("can't create transfer")
		return
	}

	newUser, err := s.userRepository.FindByID(ctx, userData.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to modify balance", err)
		err = fmt.Errorf("failed to modify balance")
		return
	}

	newUser.Balance -= transaction.Amount

	err = s.userRepository.Create(ctx, &newUser)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create new user", err)
		err = fmt.Errorf("failed to create new user")
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
