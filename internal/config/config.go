package config

import (
	"encoding/json"
	"fmt"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"log"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type AppConfig struct {
	Name  string
	Port  string
	Page  int64
	Limit int64
	Sort  string
}

type DbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	Debug           bool
}

type Config struct {
	App *AppConfig
	Db  *DbConfig
}

var config Config

func App() *AppConfig {
	return config.App
}

func Db() *DbConfig {
	return config.Db
}

func LoadConfig() {
	setDefaultConfig()

	_ = viper.BindEnv("CONSUL_URL")
	_ = viper.BindEnv("CONSUL_PATH")

	consulURL := viper.GetString("CONSUL_URL")
	consulPath := viper.GetString("CONSUL_PATH")

	if consulURL != "" && consulPath != "" {
		_ = viper.AddRemoteProvider("consul", consulURL, consulPath)

		viper.SetConfigType("yaml")
		err := viper.ReadRemoteConfig()

		if err != nil {
			log.Println(fmt.Sprintf("%s named \"%s\"", err.Error(), consulPath))
		}

		config = Config{}

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}

		if r, err := json.MarshalIndent(&config, "", "  "); err == nil {
			fmt.Println(string(r))
		}
	} else {
		logger.Info("CONSUL_URL or CONSUL_PATH missing! Serving with default config...")
	}
}

func setDefaultConfig() {
	config.App = &AppConfig{
		Name:  "go-starter",
		Port:  "8081",
		Page:  1,
		Limit: 10,
		Sort:  "created_at desc",
	}

	config.Db = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Pass:            "12345678",
		Schema:          "demo_database",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}
}
