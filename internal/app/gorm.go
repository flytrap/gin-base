package app

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/flytrap/gin-base/internal/app/config"
	"github.com/flytrap/gin-base/internal/app/repositories"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, err := NewGormDB()
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {}

	if cfg.EnableAutoMigrate {
		err = repositories.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

func NewGormDB() (*gorm.DB, error) {
	cfg := config.C.Gorm

	gConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
	}
	var dialector gorm.Dialector
	switch strings.ToLower(cfg.DBType) {
	case "mysql":
		// create database if not exists
		dsn := config.C.MySQL.DSN()
		cfgMs, err := mysqlDriver.ParseDSN(dsn)
		if err != nil {
			return nil, err
		}

		err = createDatabaseWithMySQL(cfgMs)
		if err != nil {
			return nil, err
		}

		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(config.C.Postgres.DSN())
	default:
		dialector = sqlite.Open(config.C.Sqlite3.DSN())
	}

	db, err := gorm.Open(dialector, gConfig)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	return db, nil
}

func createDatabaseWithMySQL(cfg *mysqlDriver.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", cfg.User, cfg.Passwd, cfg.Addr)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", cfg.DBName)
	_, err = db.Exec(query)
	return err
}
