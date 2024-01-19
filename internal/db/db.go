package db

import (
	"fmt"
	"my_app/internal/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	conf := config.GetC()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Env.Mysql.User, conf.Env.Mysql.Password, conf.Env.Mysql.Address, conf.Env.Mysql.Database,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to get DB instance")
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	println("initialized database")
}

func MirateTable(tables ...interface{}) {
	DB.AutoMigrate(tables...)
}
