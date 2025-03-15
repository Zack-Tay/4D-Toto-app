package main

import (
	"log"

	"github.com/Zack-Tay/4D-Toto-app/backend/toto/totoScraper"
	"github.com/gin-gonic/gin"
)

/*
WinningNumbers is a 2d array because we need to know the group the number is in
and the numbers that come out
*/
type Result4d struct {
	DrawDate       string  `json:"draw_date"`
	WinningNumbers [][]int `json:"winning_numbers"`
}

func main() {
	// Scrape the results
	r := gin.Default()

	// Endpoint to get the results
	r.GET("/api", func(c *gin.Context) {
		// Scrape results
		results, err := totoScraper.ScrapeTotoResults()
		if err != nil {
			log.Fatal("Error scraping results:", err)
			c.JSON(500, gin.H{"error": "Error scraping results"})
			return
		}

		if len(results) > 0 {
			// Returning the first result as a key-value JSON object
			result := results[0] // Assuming you want the first result
			c.JSON(200, gin.H{
				"draw_date":         result.DrawDate,
				"winning_numbers":   result.WinningNumbers,
				"additional_number": result.AdditionalNumber,
				"group1_prize_pool": result.Group1PrizePool,
				"winning_shares":    result.WinningShares,
			})
		} else {
			c.JSON(200, gin.H{"message": "No results found"})
		}
	})

	// Start the API server
	r.Run(":8080")
}
