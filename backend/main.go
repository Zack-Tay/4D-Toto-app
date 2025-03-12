package main

import (
	"fmt"
	"log"

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
	DrawDate       string `json:"draw_date"`
	WinningNumbers []int  `json:"winning_numbers"`
}

// In-memory store for results (for now)
var results = []ResultToto{
	{
		DrawDate:       "2025-03-10",
		WinningNumbers: []int{1, 2, 3, 4, 5},
	},
	{
		DrawDate:       "2025-03-09",
		WinningNumbers: []int{6, 7, 8, 9, 10},
	},
}

func ScrapeResults() ([]ResultToto, error) {
	url := "https://www.singaporepools.com.sg/en/product/pages/toto_results.aspx"
	c := colly.NewCollector()
	var results []ResultToto

	c.OnHTML(".result-row", func(e *colly.HTMLElement) {
		// Extract draw date
		drawDate := e.ChildText(".draw-date") // Adjust this selector based on the actual HTML structure

		// Extract the winning numbers
		var winningNumbers []int
		e.ForEach(".numbers .number", func(i int, el *colly.HTMLElement) {
			// Extract each number and convert it to an integer
			number := el.Text
			var num int
			fmt.Sscanf(number, "%d", &num)
			winningNumbers = append(winningNumbers, num)
		})

		// Store the result in the slice
		results = append(results, ResultToto{
			DrawDate:       drawDate,
			WinningNumbers: winningNumbers,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %s\n", r.Request.URL, err)
	})

	// Make the HTTP GET request to scrape the page
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func main() {
	// Scrape the results
	results, err := ScrapeResults()
	if err != nil {
		log.Fatal("Error scraping results: ", err)
	}

	// Print out the results
	for _, result := range results {
		fmt.Printf("Draw Date: %s, Winning Numbers: %v\n", result.DrawDate, result.WinningNumbers)
	}
}
