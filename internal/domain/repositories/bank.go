package repositories

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"gorm.io/gorm"
)

type Bank interface {
	Create(ctx context.Context, entity *entities.Bank) (err error)
	FindByUserID(ctx context.Context, userID int) (bank []entities.Bank, err error)
	FindByID(ctx context.Context, id int) (bank entities.Bank, err error)
}

type bank struct {
	db *gorm.DB
}

func NewBank(db *gorm.DB) Bank {
	if db == nil {
		panic("db is nil")
	}

	return &bank{db: db}
}

func (r *bank) Create(ctx context.Context, entity *entities.Bank) (err error) {
	err = r.db.WithContext(ctx).Create(&entity).Error
	return
}

func (r *bank) FindByUserID(ctx context.Context, userID int) (bank []entities.Bank, err error) {
	err = r.db.WithContext(ctx).Where(&entities.Bank{UserID: userID}).Find(&bank).Error
	return
}

func (r *bank) FindByID(ctx context.Context, id int) (bank entities.Bank, err error) {
	err = r.db.WithContext(ctx).Where(&entities.Bank{ID: id}).First(&bank).Error
	return
}
