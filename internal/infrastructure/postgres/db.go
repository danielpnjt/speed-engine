package postgres

import (
	"fmt"

	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/labstack/gommon/color"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDB(cfg config.PostgresqlDB) (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		panic(err)
	}

	if cfg.Debug {
		db = db.Debug()
	}

	color.Println(color.Green(fmt.Sprintf("⇨ connected to postgresql db on %s\n", cfg.Name)))
	return
}
