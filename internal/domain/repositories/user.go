package repositories

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type User interface {
	Create(ctx context.Context, entity *entities.User) (err error)
	FindByUsername(ctx context.Context, username string) (user entities.User, err error)
	FindByID(ctx context.Context, id int) (user entities.User, err error)
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.User, count int64, err error)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	if db == nil {
		panic("db is nil")
	}

	return &user{db: db}
}

func (r *user) Create(ctx context.Context, entity *entities.User) (err error) {
	err = r.db.WithContext(ctx).Save(&entity).Error
	return
}

func (r *user) FindByUsername(ctx context.Context, username string) (user entities.User, err error) {
	err = r.db.WithContext(ctx).Where(&entities.User{Username: username}).First(&user).Error
	return
}

func (r *user) FindByID(ctx context.Context, id int) (user entities.User, err error) {
	err = r.db.WithContext(ctx).Where(&entities.User{ID: id}).First(&user).Error
	return
}

func (r *user) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.User, count int64, err error) {
	limit := pagination.Limit
	offset := (pagination.Page - 1) * pagination.Limit
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		queryPayload := r.db.WithContext(egCtx).Limit(int(limit)).Offset(int(offset))
		return utils.CompileConds(queryPayload, conds...).Find(&result).Error
	})
	eg.Go(func() (egErr error) {
		countPayload := r.db.WithContext(egCtx).Model(&entities.User{})
		return utils.CompileConds(countPayload, conds...).Count(&count).Error
	})
	err = eg.Wait()
	return
}
