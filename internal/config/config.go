package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Ytcwd struct {
	MaxVideoAge            time.Duration
	NoChannelRetryInterval time.Duration
	PostFetchSleepDuration time.Duration
	MaxLastFetchAge        time.Duration
}

type General struct {
	AppEnv string
}

type Config struct {
	Logger  zerolog.Logger
	DB      DBConfig
	General General
	Ytcwd   Ytcwd
}

var once sync.Once
var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}
	_ = LoadConfig()
	return config
}

func LoadConfig() error {
	var err error

	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("./config")

		if e := viper.ReadInConfig(); e != nil {
			err = e
			return
		}

		viper.SetConfigName("config.local")
		_ = viper.MergeInConfig()

		config = &Config{
			DB: DBConfig{
				Host:     viper.GetString("database.host"),
				Port:     viper.GetInt("database.port"),
				User:     viper.GetString("database.user"),
				Password: viper.GetString("database.password"),
				Name:     viper.GetString("database.name"),
				SSLMode:  viper.GetString("database.sslmode"),
			},
			General: General{
				AppEnv: viper.GetString("general.app_env"),
			},
			Ytcwd: Ytcwd{
				MaxVideoAge:            viper.GetDuration("ytcwd.max_video_age"),
				NoChannelRetryInterval: viper.GetDuration("ytcwd.no_channel_retry_interval"),
				PostFetchSleepDuration: viper.GetDuration("ytcwd.post_fetch_sleep_duration"),
				MaxLastFetchAge:        viper.GetDuration("ytcwd.max_last_fetch_age"),
			},
		}
	})

	return err
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name, c.DB.SSLMode,
	)
}
