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
	DrawDate          string   `json:"draw_date"`
	FirstPrize        int      `json:"firstPrize"`
	SecondPrize       int      `json:"secondPrize"`
	ThirdPrize        int      `json:"thirdPrize"`
	StarterPrizes     []string `json:"startPrizes"`
	ConsolationPrizes []string `json:"consolationPrizes"`
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

func Scrape4DResults() ([]Result4d, error) {
	url := "https://www.singaporepools.com.sg/en/product/pages/4d_results.aspx?sppl=RHJhd051bWJlcj01MzAw"
	c := colly.NewCollector()
	var results []Result4d

	c.OnHTML(".four-d-results.article-body", func(e *colly.HTMLElement) {
		// drawDate
		drawDate := e.ChildText(".drawDate")
		// drawNumber := e.ChildText(".drawNumber")
		firstPrizeElem := e.ChildText(".tdFirstPrize")   // Extract text from the .tdFirstPrize element
		secondPrizeElem := e.ChildText(".tdSecondPrize") // Extract text from the .tdSecondPrize element
		thirdPrizeElem := e.ChildText(".tdThirdPrize")   // Extract text from the .tdThirdPrize element

		var firstPrize int
		var secondPrize int
		var thirdPrize int

		fmt.Sscanf(firstPrizeElem, "%d", &firstPrize)
		fmt.Sscanf(secondPrizeElem, "%d", &secondPrize)
		fmt.Sscanf(thirdPrizeElem, "%d", &thirdPrize)

		var starterPrizes []string
		var consolationPrizes []string

		// Scrape Starter Prizes
		e.ForEach(".tbodyStarterPrizes tr", func(i int, row *colly.HTMLElement) {
			row.ForEach("td", func(j int, cell *colly.HTMLElement) {
				// Parse each number in the Starter Prizes
				cellText := strings.TrimSpace(cell.Text)
				if cellText != "" {
					// Add the cleaned-up cell value to the starterPrizes array
					starterPrizes = append(starterPrizes, cellText)
				}
			})
		})

		// Scrape Consolation Prizes
		e.ForEach(".tbodyConsolationPrizes tr", func(i int, row *colly.HTMLElement) {
			row.ForEach("td", func(j int, cell *colly.HTMLElement) {
				// Parse each number in the Starter Prizes
				cellText := strings.TrimSpace(cell.Text)
				if cellText != "" {
					// Add the cleaned-up cell value to the starterPrizes array
					consolationPrizes = append(consolationPrizes, cellText)
				}
			})
		})

		// Append the extracted result to the results slice
		results = append(results, Result4d{
			DrawDate:          drawDate,
			FirstPrize:        firstPrize,
			SecondPrize:       secondPrize,
			ThirdPrize:        thirdPrize,
			StarterPrizes:     starterPrizes,
			ConsolationPrizes: consolationPrizes,
		})
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %s\n", r.Request.URL, err)
	})

	// Start scraping by visiting the 4D results URL
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ScrapeTotoResults() ([]ResultToto, error) {
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
	c := colly.NewCollector()

	url := "https://www.singaporepools.com.sg/en/product/pages/4d_results.aspx?sppl=RHJhd051bWJlcj01MzAw"

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err) // Handle any errors that occur during visiting the page
	}

	// Endpoint to get the results
	r.GET("/api/toto", func(c *gin.Context) {
		// Scrape results
		results, err := ScrapeTotoResults()
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

	// Endpoint to get the results for 4D
	r.GET("/api/4d", func(c *gin.Context) {
		// Scrape 4D results
		results, err := Scrape4DResults()
		if err != nil {
			log.Fatal("Error scraping 4D results:", err)
			c.JSON(500, gin.H{"error": "Error scraping 4D results"})
			return
		}

		if len(results) > 0 {
			// Returning the first result as a key-value JSON object
			result := results[0] // Assuming you want the first result
			c.JSON(200, gin.H{
				"draw_date":          result.DrawDate,
				"first_prize":        result.FirstPrize,
				"second_prize":       result.SecondPrize,
				"third_prize":        result.ThirdPrize,
				"starter_prizes":     result.StarterPrizes,
				"consolation_prizes": result.ConsolationPrizes,
			})
		} else {
			c.JSON(200, gin.H{"message": "No 4D results found"})
		}
	})

	// Start the API server
	r.Run(":8080")
}
