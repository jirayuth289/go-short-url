package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

func addShortUrl(c *gin.Context) {
	db := getDbConn()
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
}

func redirectShortUrl(c *gin.Context) {
	db := getDbConn()
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
}

func getShortUrlById(c *gin.Context) {
	db := getDbConn()
	fmt.Println(c.Param("id"))
	idQuery, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var shortURLs []ShortURL
	result := db.Find(&shortURLs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for _, shortURL := range shortURLs {
		if shortURL.ID == uint(idQuery) {
			c.JSON(http.StatusOK, gin.H{"shortURLs": []ShortURL{shortURL}})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"shortURLs": shortURLs})
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
