package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
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

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&ShortURL{})

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/shorten", func(c *gin.Context) {
		var data struct {
			URL string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var url ShortURL
		result := db.Where("original_url = ?", data.URL).First(&url)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				shortURL := generateShortURL()
				url = ShortURL{OriginalURL: data.URL, ShortURL: shortURL}
				result = db.Create((&url))
				if result != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"shortURL": url.ShortURL})
	})

	r.GET("/:shortURL", func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		var url ShortURL
		result := db.Where("short_url = ?", shortURL).Find(&url)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}

		c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
	})

	r.Run(":8080")
}

func generateShortURL() string {
	const chars = "abedefghijkImnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// rand.Seed(time.Now().UnixNano())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var shortURL string
	for i := 0; i < 6; i++ {
		shortURL += string(chars[r.Intn(len(chars))])
	}

	return shortURL
}
