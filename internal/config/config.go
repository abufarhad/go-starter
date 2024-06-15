package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

type AppCfg struct {
	Name  string
	Port  string
	Page  int64
	Limit int64
	Sort  string
}

type DbCfg struct {
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

type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ContextKey         string
}

type Config struct {
	App *AppCfg
	Db  *DbCfg
	Jwt *JwtConfig
}

var config Config

func App() *AppCfg {
	return config.App
}

func Jwt() *JwtConfig {
	return config.Jwt
}

func Db() *DbCfg {
	return config.Db
}

func LoadConfig() *Config {
	if _, err := os.ReadFile(".env"); err != nil {
		fmt.Errorf("env file not found %v", err)
	}

	loadConfig()

	return &config
}

func loadConfig() {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Errorf("fatal error config load: %w", err)
	}

	config.App = &AppCfg{}    // Initialize App config
	config.Db = &DbCfg{}      // Initialize Db config
	config.Jwt = &JwtConfig{} // Initialize Db config

	// App
	config.App.Port = viper.GetString("APP_PORT")

	//Database
	config.Db.Host = viper.GetString("DB_HOST")
	config.Db.Port = viper.GetString("DB_PORT")
	config.Db.User = viper.GetString("DB_USERNAME")
	config.Db.Pass = viper.GetString("DB_PASSWORD")
	config.Db.Schema = viper.GetString("DB_NAME")

	//Jwt
	config.Jwt.AccessTokenSecret = viper.GetString("ACCESS_TOKEN_SECRET")
	config.Jwt.RefreshTokenSecret = viper.GetString("REFRESH_TOKEN_SECRET")
	config.Jwt.AccessTokenExpiry = time.Minute * time.Duration(viper.GetInt("ACCESS_TOKEN_EXPIRY"))
	config.Jwt.RefreshTokenExpiry = time.Minute * time.Duration(viper.GetInt("REFRESH_TOKEN_EXPIRY"))
	config.Jwt.ContextKey = viper.GetString("CONTEXT_KEY")
}
