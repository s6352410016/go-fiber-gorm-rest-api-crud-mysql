package database

import (
	"fmt"
	"log"
	"os"

	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed To Connect Database\n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected To Database Successfully")
	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(&models.Employee{})
	DB = db
}
