package conn

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/utils/methodutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"os"
	"time"
)

var db *gorm.DB

func getDsn() string {
	var dsn = fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)
	logger.Info(dsn)
	return dsn
}

func ConnectDb() error {

	methodutil.LoadEnv()
	logger.Info("connecting to mysql at " + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "...")
	logMode := gormlogger.Info

	//dB, err := gorm.Open(postgres.Open(getDsn()), &gorm.Config{
	//	PrepareStmt: true,
	//	Logger:      gormlogger.Default.LogMode(logMode),
	//})

	//Mysql connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	fmt.Println("dsn = ", dsn)
	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormlogger.Default.LogMode(logMode),
	})

	if err != nil {
		logger.Error("mysql connection error ", err)
		return err
	}

	sqlDb, err := dB.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(0)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Hour)

	db = dB

	logger.Info("db connection successful...")
	return nil
}

func Db() *gorm.DB {
	return db
}
