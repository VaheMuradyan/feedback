package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	var err error

	// Replace with your actual database credentials
	dsn := "root:@tcp(127.0.0.1:3306)/feedbeck?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to retrieve SQL DB object:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Database connection successful!")

	// Migrate tables if necessary
	DB.AutoMigrate(&feedbeck{}, &AdminRating{})

	// Start the Gin server
	router := gin.Default()
	router.POST("/feedback", CreateFeedbackAndUpdateRating)
	router.Run(":8080")
}
