package repositories

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"gorm.io/gorm"
)

type Transaction interface {
	Create(ctx context.Context, entity *entities.Transaction) (err error)
	FindByUserID(ctx context.Context, userID int) (transaction entities.Transaction, err error)
	FindByID(ctx context.Context, id int) (transaction entities.Transaction, err error)
	FindByReference(ctx context.Context, reference string) (transaction entities.Transaction, err error)
}

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) Transaction {
	if db == nil {
		panic("db is nil")
	}

	return &transaction{db: db}
}

func (r *transaction) Create(ctx context.Context, entity *entities.Transaction) (err error) {
	err = r.db.WithContext(ctx).Create(&entity).Error
	return
}

func (r *transaction) FindByUserID(ctx context.Context, userID int) (transaction entities.Transaction, err error) {
	err = r.db.WithContext(ctx).Where(&entities.Transaction{UserID: userID}).First(&transaction).Error
	return
}

func (r *transaction) FindByID(ctx context.Context, id int) (transaction entities.Transaction, err error) {
	err = r.db.WithContext(ctx).Where(&entities.Transaction{ID: id}).First(&transaction).Error
	return
}

func (r *transaction) FindByReference(ctx context.Context, reference string) (transaction entities.Transaction, err error) {
	err = r.db.WithContext(ctx).Where(&entities.Transaction{Reference: reference}).First(&transaction).Error
	return
}
