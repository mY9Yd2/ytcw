package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"ytcw/internal/config"
)

var once sync.Once
var DB *gorm.DB

func Connect() *gorm.DB {
	if DB != nil {
		return DB
	}

	once.Do(func() {
		cfg := config.LoadConfig()
		dsn := cfg.GetDSN()

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(getLogLevel(*cfg)),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect database")
		}

		DB = db
	})

	return DB
}

func getLogLevel(cfg config.Config) logger.LogLevel {
	switch cfg.General.AppEnv {
	case "production", "prod":
		return logger.Silent
	case "debug":
		return logger.Info
	default:
		return logger.Warn
	}
}
