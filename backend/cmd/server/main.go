package main

import (
	"log"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/middleware"
	"main/internal/scraper"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialise database:", err)
	}

	router := gin.Default()
	router.Use(middleware.CORS())

	s := scraper.NewScraper()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SG Lottery API",
			"endpoints": []string{
				"GET /health",
				"GET /api/v1/results/latest",
				"POST /api/v1/scrape",
			},
		})
	})

	router.GET("/api/results/4d", func(c *gin.Context) {
		results, err := s.Scrape4DLatest()
		if err != nil {
			log.Println("Error scraping 4D results:", err)
			c.JSON(500, gin.H{"error": "Failed to scrape 4D results"})
			return
		}

		if len(results) == 0 {
			c.JSON(200, gin.H{"message": "No 4D results found"})
			return
		}

		c.JSON(200, results)
	})

	router.GET("/api/results/toto", func(c *gin.Context) {
		results, err := s.ScrapeTotoLatest()
		if err != nil {
			log.Println("Error scraping Toto results:", err)
			c.JSON(500, gin.H{"error": "Failed to scrape Toto results"})
			return
		}

		if len(results) == 0 {
			c.JSON(200, gin.H{"message": "No Toto results found"})
			return
		}

		c.JSON(200, results)
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/results/latest", handlers.GetLatestResults)
		v1.GET("/results/4d/:date", handlers.Get4DResults)
		v1.GET("/results/toto/:date", handlers.GetTotoResults)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
