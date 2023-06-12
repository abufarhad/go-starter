package config

import (
	"time"

	_ "github.com/spf13/viper/remote"
)

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
