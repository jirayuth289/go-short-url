package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type ShortURL struct {
	gorm.Model
	OriginalURL string `gorm:"unique"`
	ShortURL    string `gorm:"unique"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot loading .env file")
	}

	getDbConn()

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/shorten", addShortUrl)

	r.GET("/:shortURL", redirectShortUrl)

	r.GET("/short-url/:id", getShortUrlById)

	r.Run("localhost:8080")
}
