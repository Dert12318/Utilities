package sql

import (
	"gorm.io/gorm"
)

type DatabaseManager interface {
	GetMaster() *gorm.DB
	StartTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) *gorm.DB
	RollbackTransaction(tx *gorm.DB) *gorm.DB
}

type DatabaseConfig struct {
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	DbMaxIdleConns int
	DbMaxOpenConns int
	DbLogLevel     string
}
