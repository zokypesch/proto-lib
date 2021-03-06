package core

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// InitDB for initial database
func InitDB(address string, dbName string, user string, password string, DBPort int, logMode bool, maxOpenConn int, maxIdleConn int) *gorm.DB {
	dbUser := user
	dbPass := password
	dbEndpoint := address
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbEndpoint, DBPort, dbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	db.LogMode(logMode)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(10))
	db.DB().SetMaxIdleConns(maxIdleConn)
	db.DB().SetMaxOpenConns(maxOpenConn)
	db.SingularTable(true)
	return db
}

// InitDBWithoutMax for initial database
func InitDBWithoutMax(address string, dbName string, user string, password string, DBPort int, logMode bool) *gorm.DB {
	dbUser := user
	dbPass := password
	dbEndpoint := address
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbEndpoint, DBPort, dbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	db.LogMode(logMode)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(10))
	db.SingularTable(true)
	return db
}

// InitDBWithoutLog for initial database
func InitDBWithoutLog(address string, dbName string, user string, password string, DBPort int) *gorm.DB {
	dbUser := user
	dbPass := password
	dbEndpoint := address
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbEndpoint, DBPort, dbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(10))
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(50)
	db.SingularTable(true)
	return db
}
