package handlers

import (
	"main/internal/scraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLatestResults(c *gin.Context) {
	s := scraper.NewScraper()

	results, err := s.GetLatest4DFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"count":   len(results),
		"results": results,
	})
}

func Get4DResults(c *gin.Context) {
	date := c.Param("date")

	// TODO: Validate date format
	// TODO: Query database

	mockData := gin.H{
		"status":  "success",
		"date":    date,
		"message": "4D results endpoint - coming soon",
	}

	c.JSON(http.StatusOK, mockData)
}

func GetTotoResults(c *gin.Context) {
	date := c.Param("date")

	mockData := gin.H{
		"status":  "success",
		"date":    date,
		"message": "Toto results endpoint - coming soon",
	}

	c.JSON(http.StatusOK, mockData)
}

func ManualScrape(c *gin.Context) {
	s := scraper.NewScraper()

	if err := s.ScrapeAndSave4DResults(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Scraping completed successfully",
	})
}
