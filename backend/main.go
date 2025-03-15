package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

/*
WinningNumbers is a 2d array because we need to know the group the number is in
and the numbers that come out
*/
type Result4d struct {
	DrawDate       string  `json:"draw_date"`
	WinningNumbers [][]int `json:"winning_numbers"`
}

type ResultToto struct {
	DrawDate         string     `json:"draw_date"`
	WinningNumbers   []int      `json:"winning_numbers"`
	AdditionalNumber int        `json:"additional-number"`
	Group1PrizePool  string     `json:"group1-prize-pool"`
	WinningShares    [][]string `json:"winning-shares"`
}

func parseCurrency(currency string) (int, error) {
	// Remove the dollar sign and commas
	cleaned := strings.TrimPrefix(currency, "$")   // Remove the dollar sign
	cleaned = strings.ReplaceAll(cleaned, ",", "") // Remove the commas

	// Parse the cleaned string into an integer
	value, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}

	return value, nil
}

func ScrapeResults() ([]ResultToto, error) {
	url := "https://www.singaporepools.com.sg/en/product/sr/Pages/toto_results.aspx?sppl=RHJhd051bWJlcj00MDU5"
	c := colly.NewCollector()
	var results []ResultToto

	c.OnHTML(".toto-result.article-body", func(e *colly.HTMLElement) {
		// Extract draw date

		drawDate := e.ChildText(".drawDate") // Adjust this selector based on the actual HTML structure

		// Extract the winning numbers
		e.ForEach(".table.table-striped tbody tr", func(i int, row *colly.HTMLElement) {
			// Extract each number and convert it to an integer
			var additionalNumber int
			var prizePool string
			var winningNumbers []int
			var winningShares [][]string

			additionalNum := e.ChildText(".additional")
			prizePoolString := e.ChildText(".jackpotPrize")

			row.ForEach("td", func(j int, el *colly.HTMLElement) {
				// Extract the number from each cell
				var num int
				fmt.Sscanf(el.Text, "%d", &num)
				winningNumbers = append(winningNumbers, num)
			})

			e.ForEach(".tableWinningShares tbody tr", func(i int, row *colly.HTMLElement) {
				// Skip the header row
				if i == 0 {
					return
				}
				shareData := []string{}
				row.ForEach("td", func(j int, el *colly.HTMLElement) {
					cleanText := strings.TrimSpace(el.Text)
					shareData = append(shareData, cleanText)
				})
				winningShares = append(winningShares, shareData)
			})

			fmt.Sscanf(additionalNum, "%d", &additionalNumber)
			fmt.Sscanf(prizePool, "%s", &prizePoolString)

			results = append(results, ResultToto{
				DrawDate:         drawDate,
				WinningNumbers:   winningNumbers,
				AdditionalNumber: additionalNumber,
				Group1PrizePool:  prizePoolString,
				WinningShares:    winningShares,
			})
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %s\n", r.Request.URL, err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func main() {
	// Scrape the results
	r := gin.Default()

	// Endpoint to get the results
	r.GET("/api", func(c *gin.Context) {
		// Scrape results
		results, err := ScrapeResults()
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
