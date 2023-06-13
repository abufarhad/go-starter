package conn

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/config"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/utils/methodutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

// this func is for postgres only
func getDsn(dbCfg *config.DbCfg) string {
	var dsn = fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Port,
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Pass,
		dbCfg.Schema,
	)
	logger.Info(dsn)
	return dsn
}

func ConnectDb(dbCfg *config.DbCfg) error {
	methodutil.LoadEnv()
	logger.Info("connecting to mysql at " + dbCfg.Host + ":" + dbCfg.Port + "...")
	logMode := gormlogger.Info

	//dB, err := gorm.Open(postgres.Open(getDsn(dbCfg)), &gorm.Config{
	//	PrepareStmt: true,
	//	Logger:      gormlogger.Default.LogMode(logMode),
	//})

	//Mysql connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Schema)
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
