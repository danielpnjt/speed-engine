package queue

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/danielpnjt/speed-engine/internal/usecase/transaction"
	"gorm.io/gorm"
)

type Worker interface {
	Run()
}

type worker struct {
	transactionService transaction.Service
	machineryServer    *machinery.Server
	db                 *gorm.DB
}

func New() *worker {
	return &worker{}
}

func (w *worker) SetTransactionService(service transaction.Service) *worker {
	w.transactionService = service
	return w
}

func (w *worker) SetMachineryServer(server *machinery.Server) *worker {
	w.machineryServer = server
	return w
}

func (w *worker) SetDB(db *gorm.DB) *worker {
	w.db = db
	return w
}

func (w *worker) RegisterTasks() *worker {
	if w.transactionService == nil {
		panic("transactionService is nil")
	}

	if w.db == nil {
		panic("db is nil")
	}

	w.machineryServer.RegisterTasks(map[string]interface{}{
		"enqueue-speed_engine-topup": w.EnqueueStatusTopUp,
	})

	return w
}

func (w *worker) Run() {
	worker := w.machineryServer.NewWorker("speed_engine-worker", config.GetInt("worker.totalWorker"))
	err := worker.Launch()
	if err != nil {
		panic(err)
	}
}
