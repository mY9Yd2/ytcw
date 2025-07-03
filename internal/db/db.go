package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect database")
		}

		DB = db
	})

	return DB
}
