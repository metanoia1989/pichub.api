package database

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"pichub.api/config"
)

var (
	DB  *gorm.DB
	err error
)

// DbConnection create database connection
func DbConnection() error {
	var db = DB

	logMode := viper.GetBool("DB_LOG_MODE")

	loglevel := logger.Silent
	if logMode {
		loglevel = logger.Info
	}

	var dbConfig = config.Config.Database

	// log.Printf("dbConfig %+v", dbConfig)
	db, err = gorm.Open(mysql.Open(dbConfig.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
		Logger:                                   logger.Default.LogMode(loglevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: dbConfig.Prefix,
		},
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
