package models

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	err := godotenv.Load("/etc/mysql-app.conf")
	if err != nil {
		slog.Error("godotenv load", "error", err)
		os.Exit(1)
	}
	username := os.Getenv("db_user_name")
	password := os.Getenv("db_user_pass")
	dbName := os.Getenv("db_app_name")
	dbHost := os.Getenv("db_app_host")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbHost, dbName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("open db", "error", err)
	}
	db = conn
	db.AutoMigrate(&Location{})

}

func GetDB() *gorm.DB {
	return db
}
