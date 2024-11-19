package handler

import (
	"github.com/danielpnjt/speed-engine/internal/infrastructure/container"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	SpeedEngineDB      *gorm.DB
	healthCheckHandler *healthCheckHandler
	userHandler        *userHandler
	bankHandler        *bankHandler
	transactionHandler *transactionHandler
	adminHandler       *adminHandler
	redisClient        *redis.Client
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		SpeedEngineDB:      container.SpeedEngineDB,
		healthCheckHandler: NewHealthCheckHandler().SetHealthCheckService(container.HealthCheckService).Validate(),
		userHandler:        NewUserHandler().SetUserService(container.UserService).Validate(),
		bankHandler:        NewBankHandler().SetBankService(container.BankService).Validate(),
		transactionHandler: NewTransactionHandler().SetTransactionService(container.TransactionService).Validate(),
		redisClient:        container.RedisClient,
	}
}

func (h *Handler) Validate() *Handler {
	if h.SpeedEngineDB == nil {
		panic("SpeedEngineDB is nil")
	}
	if h.healthCheckHandler == nil {
		panic("healthCheckHandler is nil")
	}
	if h.userHandler == nil {
		panic("userHandler is nil")
	}
	if h.bankHandler == nil {
		panic("bankHandler is nil")
	}
	if h.transactionHandler == nil {
		panic("transactionHandler is nil")
	}
	if h.redisClient == nil {
		panic("redisClient is nil")
	}
	return h
}
