package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Dert12318/Utilities/database/sql"
)

type PostgresDatabaseManager struct {
	Master *gorm.DB
}

func NewPostgresDatabase(conf sql.DatabaseConfig) sql.DatabaseManager {
	dsn := PostgresURI(
		conf.DBUser,
		conf.DBPassword,
		fmt.Sprintf(`%s:%s`, conf.DBHost, conf.DBPort),
		conf.DBName)

	logLevel := logger.Warn // warn is default gorm log level
	switch strings.ToLower(conf.DbLogLevel) {
	case "info":
		logLevel = logger.Info
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}

	db, err := Initialize(dsn, conf.DbMaxIdleConns, conf.DbMaxOpenConns, logLevel)
	if err != nil {
		log.Fatal("error")
	}

	database := PostgresDatabaseManager{
		Master: db,
	}
	return &database
}

func Initialize(dsn string, maxIdleConns int, maxOpenConns int, logLevel logger.LogLevel) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 2)

	return db, nil
}

func (dbManager *PostgresDatabaseManager) GetMaster() *gorm.DB {
	if dbManager.Master == nil {
		return nil
	}
	return dbManager.Master
}

func (dbManager *PostgresDatabaseManager) StartTransaction() *gorm.DB {
	return dbManager.Master.Begin()
}

func (dbManager *PostgresDatabaseManager) CommitTransaction(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (dbManager *PostgresDatabaseManager) RollbackTransaction(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}

func PostgresURI(dbUserName, dbPassword, dbAddress, dbName string) string {
	return fmt.Sprintf(`postgres://%s:%s@%s/%s?sslmode=disable`,
		dbUserName, dbPassword, dbAddress, dbName)
}
