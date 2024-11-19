package user

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	mocksRepo "github.com/danielpnjt/speed-engine/internal/domain/repositories/mocks"
	mocksWrapperRedis "github.com/danielpnjt/speed-engine/internal/infrastructure/redis/mocks"
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	token = "valid-token"
	now   = time.Now()
)

func init() {
	// log.New()
}

func TestValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// * mock gorm
	db, _, _ := sqlmock.New()
	mockGorm, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	mockRepoUser := mocksRepo.NewMockUser(ctrl)
	mockWrapperRedis := mocksWrapperRedis.NewMockWrapper(ctrl)

	service := NewService().
		SetDB(nil).
		SetUserRepository(mockRepoUser).
		SetRedisWrapper(mockWrapperRedis)

	t.Run("panic when db is nil", func(t *testing.T) {
		require.Panics(t, func() {
			service.Validate()
		}, "db is nil")
	})

	service.SetDB(mockGorm)

	t.Run("panic when userRepository is nil", func(t *testing.T) {
		service.SetUserRepository(nil)
		require.Panics(t, func() {
			service.Validate()
		}, "userRepository is nil")
	})

	service.SetUserRepository(mockRepoUser)

	t.Run("panic when redisWrapper is nil", func(t *testing.T) {
		service.SetRedisWrapper(nil)
		require.Panics(t, func() {
			service.Validate()
		}, "redisWrapper is nil")
	})

	service.SetRedisWrapper(mockWrapperRedis)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// * mock gorm
	db, _, _ := sqlmock.New()
	mockDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	mockUserRepo := mocksRepo.NewMockUser(ctrl)
	mockWrapperRedis := mocksWrapperRedis.NewMockWrapper(ctrl)

	service := &service{
		userRepository: mockUserRepo,
		redisWrapper:   mockWrapperRedis,
		db:             mockDB,
	}

	tests := []struct {
		name               string
		req                LoginRequest
		doMockUserRepo     func(mock *mocksRepo.MockUser)
		doMockRedisWrapper func(mock *mocksWrapperRedis.MockWrapper)
		wantRes            constants.DefaultResponse
		wantErr            error
	}{
		{
			name: "positive case :)",
			req: LoginRequest{
				Username: "daniel.pnjt",
				Password: "DK!@Password123",
			},
			doMockUserRepo: func(mock *mocksRepo.MockUser) {
				mock.EXPECT().FindByUsername(gomock.Any(), "daniel.pnjt").Return(
					entities.User{
						ID:       123,
						Username: "daniel.pnjt",
						Password: "$2a$10$0AiqFrWfcar1Cuq9jL8fE.XhHjkrq6O5d26Hwt1t2McO2033hcrZq",
						Email:    "daniel.pnjt@gmail.com",
						Name:     "Daniel Alexander",
						Balance:  0,
					}, nil).Times(1)
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			doMockRedisWrapper: func(mock *mocksWrapperRedis.MockWrapper) {
				mock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},

			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS,
				Data: LoginResponseData{
					Token:    token,
					ExpireAt: 1234567890,
				},
				Errors: make([]string, 0),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doMockUserRepo(mockUserRepo)
			tt.doMockRedisWrapper(mockWrapperRedis)

			_, err := service.Login(context.TODO(), tt.req)
			// require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
