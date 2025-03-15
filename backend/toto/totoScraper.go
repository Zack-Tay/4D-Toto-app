package totoScraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ResultToto stores the extracted results from the page
type ResultToto struct {
	DrawDate         string     `json:"draw_date"`
	WinningNumbers   []int      `json:"winning_numbers"`
	AdditionalNumber int        `json:"additional_number"`
	Group1PrizePool  string     `json:"group1_prize_pool"`
	WinningShares    [][]string `json:"winning_shares"`
}

// parseCurrency parses a currency string (e.g., "$1,337,645") into an integer
func parseCurrency(currency string) (int, error) {
	cleaned := strings.TrimPrefix(currency, "$")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	value, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// ScrapeTotoResults scrapes the results from a specific URL
func ScrapeTotoResults(url string) ([]ResultToto, error) {
	c := colly.NewCollector()
	var results []ResultToto

	// Scrape the page for specific data
	c.OnHTML(".divSingleDraw", func(e *colly.HTMLElement) {
		// Extract draw date
		drawDate := e.ChildText(".drawDate")

		// Extract the additional number
		additionalNum := e.ChildText(".additional")
		var additionalNumber int
		if additionalNum != "" {
			fmt.Sscanf(additionalNum, "%d", &additionalNumber)
		}

		// Extract prize pool
		prizePoolString := e.ChildText(".jackpotPrize")
		group1PrizePool, err := parseCurrency(prizePoolString)
		if err != nil {
			log.Println("Error parsing group 1 prize pool:", err)
			return
		}

		// Extract winning numbers
		var winningNumbers []int
		e.ForEach(".table.table-striped tbody tr", func(i int, row *colly.HTMLElement) {
			row.ForEach("td", func(j int, el *colly.HTMLElement) {
				var num int
				fmt.Sscanf(el.Text, "%d", &num)
				winningNumbers = append(winningNumbers, num)
			})
		})

		// Extract winning shares
		var winningShares [][]string
		e.ForEach(".tableWinningShares tbody tr", func(i int, row *colly.HTMLElement) {
			if i == 0 { // Skip header row
				return
			}
			shareData := []string{}
			row.ForEach("td", func(j int, el *colly.HTMLElement) {
				cleanText := strings.TrimSpace(el.Text)
				shareData = append(shareData, cleanText)
			})
			winningShares = append(winningShares, shareData)
		})

		// Append the scraped result
		results = append(results, ResultToto{
			DrawDate:         drawDate,
			WinningNumbers:   winningNumbers,
			AdditionalNumber: additionalNumber,
			Group1PrizePool:  fmt.Sprintf("%d", group1PrizePool),
			WinningShares:    winningShares,
		})
	})

	// Error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %s\n", r.Request.URL, err)
	})

	// Visit the URL and start scraping
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return results, nil
}
