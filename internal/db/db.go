package db

import (
	"github.com/mY9Yd2/ytcw/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var once sync.Once
var db *gorm.DB

func Connect() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error

	once.Do(func() {
		cfg := config.GetConfig()
		dsn := cfg.GetDSN()

		database, e := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(getLogLevel(*cfg)),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
		if e != nil {
			err = e
			return
		}

		db = database
	})

	return db, err
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
