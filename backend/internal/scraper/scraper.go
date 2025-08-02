package scraper

import (
	"fmt"
	"log"
	"main/internal/database"
	"main/internal/models"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

type Scraper struct {
	collector *colly.Collector
	db        *gorm.DB
}

func NewScraper() *Scraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("www.singaporepools.com.sg", "singaporepools.com.sg"),
		colly.AllowURLRevisit(),
	)

	// Add delay between requests to be respectful
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*singaporepools.com.sg*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping %s: %v", r.Request.URL, err)
	})

	return &Scraper{
		collector: c,
		db:        database.GetDB(),
	}
}

func (s *Scraper) Scrape4DLatest() ([]models.Result4DReponse, error) {
	url := "https://www.singaporepools.com.sg/DataFileArchive/Lottery/Output/fourd_result_top_draws_en.html"

	// var drawResult models.DrawResult
	var results []models.Result4DReponse

	c := s.collector.Clone()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		count := 0

		e.ForEach("li", func(i int, li *colly.HTMLElement) {
			if count >= 5 {
				return
			}
			count++

			drawDate := li.ChildText("th.drawDate")
			drawNumber := li.ChildText("th.drawNumber")
			firstPrize := li.ChildText("td.tdFirstPrize")
			secondPrize := li.ChildText("td.tdSecondPrize")
			thirdPrize := li.ChildText("td.tdThirdPrize")

			var starterPrizers []string
			li.ForEach(".tbodyStarterPrizes td", func(_ int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				if text != "" {
					starterPrizers = append(starterPrizers, text)
				}
			})

			var consolationPrizes []string
			li.ForEach(".tbodyConsolationPrizes td", func(_ int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				if text != "" {
					consolationPrizes = append(consolationPrizes, text)
				}
			})

			// build the 4d results object for that draw.
			results = append(results, models.Result4DReponse{
				DrawDate:     drawDate,
				DrawNumber:   drawNumber,
				First:        firstPrize,
				Second:       secondPrize,
				Third:        thirdPrize,
				Starters:     starterPrizers,
				Consolations: consolationPrizes,
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit 4D page: %w", err)
	}

	c.Wait()

	log.Printf("Scraped %d 4D results", len(results))
	return results, nil
}

func (s *Scraper) ScrapeAndSave4DResults() error {
	log.Println("Starting 4D scraping and saving...")

	results, err := s.Scrape4DLatest()
	if err != nil {
		return fmt.Errorf("failed to scrape: %w", err)
	}

	savedCount := 0
	for _, result := range results {
		var existing models.Result4DReponse
		err := s.db.Where("draw_number = ?", result.DrawNumber).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Doesn't exist, save it
			if err := s.db.Create(&result).Error; err != nil {
				log.Printf("Error saving draw %s: %v", result.DrawNumber, err)
			} else {
				savedCount++
				log.Printf("Saved new draw: %s", result.DrawNumber)
			}
		}
	}

	log.Printf("Saved %d new results out of %d scraped", savedCount, len(results))
	return nil
}

func (s *Scraper) GetLatest4DFromDB() ([]models.Result4DReponse, error) {
	var results []models.Result4DReponse
	err := s.db.Order("id DESC").Limit(3).Find(&results).Error
	return results, err
}

func (s *Scraper) ScrapeTotoLatest() ([]models.ResultTotoResponse, error) {
	url := "https://www.singaporepools.com.sg/DataFileArchive/Lottery/Output/toto_result_top_draws_en.html"

	var results []models.ResultTotoResponse

	c := s.collector.Clone()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		count := 0

		e.ForEach("li", func(i int, li *colly.HTMLElement) {
			if count >= 5 {
				return
			}
			count++

			drawDate := li.ChildText("th.drawDate")
			drawNumber := li.ChildText("th.drawNumber")

			var winningNums []string
			winningNumCount := 0
			li.ForEach("td", func(j int, el *colly.HTMLElement) {
				if winningNumCount >= 6 {
					return
				}

				text := strings.TrimSpace(el.Text)
				if text != "" {
					winningNums = append(winningNums, text)
				}
				winningNumCount++
			})

			additionalNumber := li.ChildText("td.additional")
			group1Prize := li.ChildText("td.jackpotPrize")

			var prizeBreakdown []models.TotoPrizeBreakdown
			li.DOM.Find("table.table.table-striped.tableWinningShares tbody tr").Each(func(i int, tr *goquery.Selection) {
				if i == 0 {
					return
				}

				var groupNum string
				var shareAmount string
				var winningShares string
				tr.Find("td").Each(func(j int, td *goquery.Selection) {
					text := strings.TrimSpace(td.Text())
					switch j {
					case 0:
						groupNum = text
					case 1:
						shareAmount = text
					case 2:
						winningShares = text
					}
				})

				prizeBreakdown = append(prizeBreakdown, models.TotoPrizeBreakdown{
					Group:         groupNum,
					ShareAmount:   shareAmount,
					WinningShares: winningShares,
				})
			})

			results = append(results, models.ResultTotoResponse{
				DrawDate:         drawDate,
				DrawNumber:       drawNumber,
				WinningNumbers:   winningNums,
				AdditionalNumber: additionalNumber,
				Group1Prize:      group1Prize,
				PrizeBreakdown:   prizeBreakdown,
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit Toto page: %w", err)
	}

	c.Wait()

	log.Printf("Scraped %d Toto results", len(results))
	return results, nil
}

// func (s *Scraper) saveResultIfNew(scraped models.Result4DReponse) (bool, error) {
// 	drawNum, err := strconv.Atoi(scraped.DrawNumber)
// 	if err != nil {
// 		return false, fmt.Errorf("invalid draw number: %s", scraped.DrawNumber)
// 	}

// 	drawDate, err := parseDrawDate(scraped.DrawDate)
// 	if err != nil {
// 		return false, fmt.Errorf("invalid draw date: %s", scraped.DrawDate)
// 	}

// 	tx := s.db.Begin()
// 	var existing models.DrawResult
// 	err = tx.Where("draw_type = ? AND draw_number = ?", "4D", drawNum).First(&existing).Error

// 	if err == nil {
// 		log.Printf("Draw %d already in database, skipping", drawNum)
// 		tx.Rollback()
// 		return false, nil
// 	}

// 	if err != gorm.ErrRecordNotFound {
// 		tx.Rollback()
// 		return false, fmt.Errorf("database error: %w", err)
// 	}

// 	if err := tx.Create(&scraped).Error; err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("failed to create draw: %w", err)
// 	}

// 	if err := tx.Create(&scraped).Error; err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("failed to create results: %w", err)
// 	}

// 	// Commit transaction
// 	if err := tx.Commit().Error; err != nil {
// 		return false, fmt.Errorf("failed to commit: %w", err)
// 	}

// 	log.Printf("Successfully saved draw %d from %s", drawNum, drawDate.Format("2006-01-02"))
// 	return true, nil
// }

// func parseDrawDate(dateStr string) (time.Time, error) {
// 	dateStr = strings.TrimSpace(dateStr)

// 	t, err := time.Parse("Mon, 2 Jan 2006", dateStr)
// 	if err == nil {
// 		return t, nil
// 	}

// 	return time.Time{}, fmt.Errorf("cannot parse date: %s", dateStr)
// }
