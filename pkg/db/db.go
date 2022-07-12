package db

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foxdex/ftx-site/pkg/consts"

	"github.com/foxdex/ftx-site/pkg/log"
	"gorm.io/gorm/schema"

	"github.com/foxdex/ftx-site/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type gormLogger struct{}

func (*gormLogger) Printf(format string, v ...interface{}) {
	format = strings.Replace(format, "\n", " ", 1)
	log.Sugar.Infof(format, v...)
}

func getLoggerLevel() logger.LogLevel {
	if os.Getenv(consts.UnitTestEnv) == "true" || gin.Mode() != gin.ReleaseMode {
		return logger.Info
	}

	return logger.Warn
}

// NewMysql connect to mysql
func NewMysql() {
	var err error
	_mysql := config.GetConfig().Mysql

	databaseURL := _mysql.Url
	newLogger := logger.New(
		&gormLogger{},
		logger.Config{
			SlowThreshold:             time.Second * time.Duration(_mysql.SlowThreshold), // 慢 SQL 阈值
			LogLevel:                  getLoggerLevel(),                                  // 日志级别
			IgnoreRecordNotFoundError: false,                                             // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                             // 禁用彩色打印
		},
	)

	mysqlConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   _mysql.Prefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: true,          // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: newLogger,
	}
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       databaseURL, // data source name
		DefaultStringSize:         255,         // default size for string fields
		DisableDatetimePrecision:  true,        // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,        // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,        // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,       // auto configure based on currently MySQL version
	}), &mysqlConfig)
	if err != nil {
		panic(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// 设置与数据库建立连接的最大数目
	mysqlDB.SetMaxOpenConns(_mysql.MaxOpenConns)
	// 设置连接池中的最大闲置连接数
	mysqlDB.SetMaxIdleConns(_mysql.MaxIdleConns)
}

// Mysql get a connection for mysql
func Mysql() *gorm.DB {
	return db
}

// DisconnectMysql disconnect mysql
func DisconnectMysql() error {
	mysqlDB, _ := db.DB()
	return mysqlDB.Close()
}
