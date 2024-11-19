package handler

import (
	"github.com/danielpnjt/speed-engine/internal/infrastructure/container"
	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, cnt *container.Container) {
	h := SetupHandler(cnt).Validate()

	e.GET("/", h.healthCheckHandler.HealthCheck)

	v1 := e.Group("/v1")
	{
		// ======== ADMIN ========
		admin := v1.Group("/admin", h.BasicAuth("admin"))
		{
			user := admin.Group("/user")
			{
				user.GET("", h.userHandler.GetAll)
				user.GET("/:userID", h.userHandler.GetDetail)
			}
		}

		// ======== PLAYER ========
		player := v1.Group("/player")
		{
			onboard := player.Group("/onboard")
			{
				onboard.POST("/register", h.userHandler.Register)
				onboard.POST("/login", h.userHandler.Login)
			}
		}

		player.Use(h.Authentication)
		{
			user := player.Group("/user")
			{
				user.GET("", h.userHandler.GetDetailPlayer)
				user.POST("/logout", h.userHandler.Logout)
			}
			bank := player.Group("/bank")
			{
				bank.GET("", h.bankHandler.FindAll)
				bank.POST("/submit-bank", h.bankHandler.SubmitBank)
			}
			transaction := player.Group("/transaction")
			{
				transaction.POST("/generate", h.transactionHandler.Generate)
				transaction.POST("/withdraw", h.transactionHandler.Withdraw)
			}
		}
	}
}
