package user

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"github.com/danielpnjt/speed-engine/internal/domain/repositories"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
	"github.com/danielpnjt/speed-engine/internal/pkg/types"
	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"gorm.io/gorm"
)

type service struct {
	db             *gorm.DB
	userRepository repositories.User
	redisWrapper   redis.Wrapper
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *gorm.DB) *service {
	s.db = db
	return s
}

func (s *service) SetUserRepository(repository repositories.User) *service {
	s.userRepository = repository
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
	if s.userRepository == nil {
		panic("userRepository is nil")
	}
	if s.redisWrapper == nil {
		panic("redisWrapper is nil")
	}
	return s
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (res constants.DefaultResponse, err error) {
	if !utils.IsValidPassword(req.Password) {
		slog.ErrorContext(ctx, "password is not valid")
		err = fmt.Errorf("password is not valid")
		return
	}

	if req.Password != req.ConfirmPassword {
		slog.ErrorContext(ctx, "password is not match")
		err = fmt.Errorf("password is not match")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		slog.ErrorContext(ctx, "failed to compile or hash new password")
		err = fmt.Errorf("failed to compile or hash new password")
		return
	}

	_, err = s.userRepository.FindByUsername(ctx, req.Username)
	if err == nil {
		slog.ErrorContext(ctx, "user already exists")
		err = fmt.Errorf("user already exists")
		return
	}

	newUser := entities.User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
		Balance:  0,
	}

	err = s.userRepository.Create(ctx, &newUser)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create new user")
		err = fmt.Errorf("failed to create new user")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    nil,
		Errors:  make([]string, 0),
	}
	return
}

func (s *service) Login(ctx context.Context, req LoginRequest) (res constants.DefaultResponse, err error) {
	var token string
	var exp int64

	user, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find user by email")
		err = fmt.Errorf("email or password is wrong")
		return
	}

	token, exp, err = utils.JwtSign(user.ID, req.Username)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to sign JWT token")
		err = fmt.Errorf("Failed to sign JWT token")
		return
	}

	newLogin := entities.Login{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Name:      user.Name,
		Balance:   user.Balance,
		Token:     token,
		ExpireAt:  exp,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	err = s.redisWrapper.Set(ctx, req.Username, time.Duration(exp)*time.Second, newLogin)
	if err != nil {
		slog.ErrorContext(ctx, "failed to store session in Redis")
		err = fmt.Errorf("failed to create session")
		return
	}

	respData := LoginResponseData{
		Token:    token,
		ExpireAt: exp,
	}
	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    respData,
		Errors:  make([]string, 0),
	}

	return
}

func (s *service) Logout(ctx context.Context) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	err = s.redisWrapper.Delete(ctx, userData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete session token from Redis")
		err = fmt.Errorf("failed to logout")
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

func (s *service) GetDetail(ctx context.Context, userID int) (res constants.DefaultResponse, err error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find user by id")
		err = fmt.Errorf("failed to find user by id")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    user,
		Errors:  make([]string, 0),
	}
	return
}

func (s *service) GetAll(ctx context.Context, req FindAllRequest) (res constants.DefaultResponse, err error) {
	users, count, err := s.userRepository.FindAllAndCount(ctx, req.PaginationRequest)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find alll")
		err = fmt.Errorf("failed to find all")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: users,
			PaginationData: constants.PaginationData{
				Page:        uint(req.Page),
				Limit:       uint(req.Limit),
				TotalPages:  uint(math.Ceil(float64(count) / float64(req.Limit))),
				TotalItems:  uint(count),
				HasNext:     req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
				HasPrevious: req.Page > 1,
			},
		},
		Errors: make([]string, 0),
	}
	return
}

func (s *service) GetDetailPlayer(ctx context.Context) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(types.String("user")).(entities.Login)
	user, err := s.userRepository.FindByID(ctx, userData.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find user by id")
		err = fmt.Errorf("failed to find user by id")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    user,
		Errors:  make([]string, 0),
	}
	return
}
