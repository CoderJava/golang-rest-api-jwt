package config

import (
	"fmt"
	"golang-rest-api-jwt/entity"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Buat koneksi ke database.
func SetupDatabaseConnection() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load env file. Make sure .env file is exists!")
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to create a connection to database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})
	fmt.Println("Database connected!")
	return db
}

// Tutup koneksi dari database.
func CloseDatabaseConnection(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	sqlDb.Close()
}
