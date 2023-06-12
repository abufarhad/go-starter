package config

import (
	"time"

	_ "github.com/spf13/viper/remote"
)

type AppConfig struct {
	App struct {
		Name  string
		Port  string
		Page  int64
		Limit int64
		Sort  string
	}

	Db struct {
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
}

type App struct {
	Name  string
	Port  string
	Page  int64
	Limit int64
	Sort  string
}

type Db struct {
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
