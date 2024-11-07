package database

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

// DbConnection create database connection
func DbConnection(masterDSN string) error {
	var db = DB

	logMode := viper.GetBool("DB_LOG_MODE")

	loglevel := logger.Silent
	if logMode {
		loglevel = logger.Info
	}

	db, err = gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})

	if err != nil {
		log.Fatalf("Db connection error")
		return err
	}
	DB = db
	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}
