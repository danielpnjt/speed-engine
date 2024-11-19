package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/danielpnjt/speed-engine/internal/domain/entities"
	"github.com/danielpnjt/speed-engine/internal/infrastructure/redis"
	"github.com/danielpnjt/speed-engine/internal/pkg/types"
	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) BasicAuth(prefix string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			basicToken := c.Request().Header.Get("Authorization")
			if basicToken == "" {
				c.Set("unauthorized", true)
				err := fmt.Errorf("authorization header is empty")
				slog.ErrorContext(ctx, "authentication failed", err)
				return err
			}
			sliceToken := strings.Split(basicToken, "Basic ")
			if len(sliceToken) < 2 {
				c.Set("unauthorized", true)
				err := fmt.Errorf("basic token slice is not greater than 2")
				slog.ErrorContext(ctx, "authentication failed", err)
				return err
			}
			token := sliceToken[1]
			username := config.GetString(fmt.Sprintf("basicAuth.%s.username", prefix))
			password := config.GetString(fmt.Sprintf("basicAuth.%s.password", prefix))
			if authToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))); authToken != token {
				c.Set("forbidden", true)
				err := fmt.Errorf("invalid basic auth provided")
				slog.ErrorContext(ctx, "authentication failed", prefix, err, token, authToken)
				return err
			}
			return next(c)
		}
	}
}

func (h *Handler) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		header := c.Request().Header
		var claims utils.JWTClaimsData
		var token string

		bearerToken := header.Get("Authorization")
		if bearerToken == "" {
			c.Set("unauthorized", true)
			err := fmt.Errorf("unauthorized [00]")
			slog.ErrorContext(ctx, "authentication failed", err)
			return err
		}
		sliceToken := strings.Split(bearerToken, "Bearer ")
		if len(sliceToken) < 2 {
			c.Set("unauthorized", true)
			err := fmt.Errorf("unauthorized [01]")
			slog.ErrorContext(ctx, "authentication failed", err)
			return err
		}
		token = sliceToken[1]
		cl, err := utils.JwtVerify(token)
		if err != nil {
			c.Set("unauthorized", true)
			slog.ErrorContext(ctx, "authentication failed", "failed to verify token", err)
			err = fmt.Errorf("unauthorized [02]")
			return err
		}
		claims.ID = cl.ID
		claims.Username = cl.Username

		userData, err := redis.NewRedisConnection(h.redisClient).Get(ctx, cl.Username)
		if err != nil {
			c.Set("unauthorized", true)
			slog.ErrorContext(ctx, "authentication failed", "failed to find user by token", err)
			err = fmt.Errorf("unauthorized [03]")
			return err
		}

		var loginRequest entities.Login
		jsonData, err := json.Marshal(userData)
		if err != nil {
			c.Set("unauthorized", true)
			slog.ErrorContext(ctx, "authentication failed", "failed marshal", err)
			err = fmt.Errorf("unauthorized [04]")
			return err
		}
		err = json.Unmarshal(jsonData, &loginRequest)
		if err != nil {
			c.Set("unauthorized", true)
			slog.ErrorContext(ctx, "authentication failed", "failed to unmarshal", err)
			err = fmt.Errorf("unauthorized [05]")
			return err
		}

		if token != loginRequest.Token {
			c.Set("unauthorized", true)
			slog.ErrorContext(ctx, "authentication failed", "token not match", err)
			err = fmt.Errorf("unauthorized [06]")
			return err
		}

		ctx = context.WithValue(ctx, types.String("user"), loginRequest)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
