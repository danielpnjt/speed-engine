package container

import (
	"fmt"
	"os"

	"github.com/RichardKnop/machinery/v1"
	machineryConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/danielpnjt/speed-engine/internal/domain/repositories"
	paymentWrap "github.com/danielpnjt/speed-engine/internal/infrastructure/payment"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/postgres"
	redisWrap "github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/worker/queue"
	"github.com/danielpnjt/speed-engine/internal/usecase/bank"
	"github.com/danielpnjt/speed-engine/internal/usecase/healthcheck"
	"github.com/danielpnjt/speed-engine/internal/usecase/transaction"
	"github.com/danielpnjt/speed-engine/internal/usecase/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	Config             *config.DefaultConfig
	HealthCheckService healthcheck.Service
	SpeedEngineDB      *gorm.DB
	UserService        user.Service
	BankService        bank.Service
	TransactionService transaction.Service
	RedisClient        *redis.Client
	QueueWorker        queue.Worker
}

func (c *Container) Validate() *Container {
	// Custom validation logic, if needed
	return c
}

func loadConfig() *config.DefaultConfig {
	config.Load(os.Getenv("env"), ".env")
	return &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("appName"),
			Address:  config.GetString("address"),
			HttpPort: config.GetString("port"),
		},
	}
}

func New() *Container {
	defConfig := loadConfig()

	redisURL := "redis://"
	if config.GetString("redis.password") != "" {
		redisURL = fmt.Sprintf("%s%s@", redisURL, config.GetString("redis.password"))
	}
	redisURL = fmt.Sprintf("%s%s", redisURL, config.GetString("redis.address"))

	workerConfig := &machineryConfig.Config{
		DefaultQueue:    "speed_engine-queue",
		ResultsExpireIn: config.GetInt("worker.taskExpiredInSecond"),
		Broker:          redisURL,
		ResultBackend:   redisURL,
		Redis: &machineryConfig.RedisConfig{
			MaxIdle:   config.GetInt("worker.maxIdle"),
			MaxActive: config.GetInt("worker.maxActive"),
		},
	}

	workerServer, err := machinery.NewServer(workerConfig)
	if err != nil {
		panic(err)
	}

	postgresqlDB := &config.PostgresqlDB{
		Host:     config.GetString("postgresql.speed_engine.host"),
		User:     config.GetString("postgresql.speed_engine.user"),
		Password: config.GetString("postgresql.speed_engine.password"),
		Name:     config.GetString("postgresql.speed_engine.db"),
		Port:     config.GetInt("postgresql.speed_engine.port"),
		SSLMode:  config.GetString("postgresql.speed_engine.ssl"),
		Schema:   config.GetString("postgresql.speed_engine.schema"),
		Debug:    config.GetBool("postgresql.speed_engine.debug"),
	}
	speedEngineDB := postgres.NewDB(*postgresqlDB)

	redisConfig := &config.RedisConfig{
		Host:     config.GetString("redis.speed_engine.host"),
		User:     config.GetString("redis.speed_engine.user"),
		Password: config.GetString("redis.speed_engine.password"),
		DB:       config.GetInt("redis.speed_engine.db"),
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host,
		Username: redisConfig.User,
		Password: redisConfig.Password,
	})
	redisWrapper := redisWrap.NewRedisConnection(redisClient)
	paymentWrapper := paymentWrap.NewPaymentWrapper()

	userRepository := repositories.NewUser(speedEngineDB)
	bankRepository := repositories.NewBank(speedEngineDB)
	transactionRepository := repositories.NewTransaction(speedEngineDB)

	healthCheckService := healthcheck.NewService().Validate()
	userService := user.NewService().
		SetDB(speedEngineDB).
		SetUserRepository(userRepository).
		SetRedisWrapper(redisWrapper).
		Validate()

	bankService := bank.NewService().
		SetDB(speedEngineDB).
		SetBankRepository(bankRepository).
		SetRedisWrapper(redisWrapper).
		Validate()

	transactionService := transaction.NewService().
		SetDB(speedEngineDB).
		SetTransactionRepository(transactionRepository).
		SetUserRepository(userRepository).
		SetBankRepository(bankRepository).
		SetRedisWrapper(redisWrapper).
		SetPaymentWrapper(paymentWrapper).
		SetWorker(workerServer).
		Validate()

	queueWorker := queue.New().
		SetMachineryServer(workerServer).
		SetTransactionService(transactionService).
		SetDB(speedEngineDB).
		RegisterTasks()

	container := &Container{
		Config:             defConfig,
		HealthCheckService: healthCheckService,
		SpeedEngineDB:      speedEngineDB,
		UserService:        userService,
		BankService:        bankService,
		TransactionService: transactionService,
		RedisClient:        redisClient,
		QueueWorker:        queueWorker,
	}
	return container.Validate()
}
