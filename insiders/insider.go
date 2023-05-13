package insiders

import (
	"fmt"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type InsiderTableHeader struct {
	FilingDate  string
	Ticker      string
	CompanyName string
	InsiderName string
	Title       string
	Price       string
	Value       string
}

func (ith *InsiderTableHeader) parseData() InsiderTableHeader {
	return InsiderTableHeader{
		FilingDate:  ith.FilingDate,
		Ticker:      ith.Ticker,
		CompanyName: ith.CompanyName,
		InsiderName: ith.InsiderName,
		Title:       ith.Title,
		Price:       ith.Price,
		Value:       ith.Value,
	}
}

func Crawl(db *gorm.DB) error {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("table.tinytable > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			var result []string
			var ith InsiderTableHeader
			tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
				result = append(result, td.Text)
			})
			ith.FilingDate = result[1]
			ith.Ticker = result[3]
			ith.CompanyName = result[4]
			ith.InsiderName = result[5]
			ith.Title = result[6]
			ith.Price = result[8]
			ith.Value = result[12]
			if isToday(ith.FilingDate) {
				db.Create(&ith)
			}
		})
	})

	err := c.Visit("http://openinsider.com/latest-insider-sales-100k")
	if err != nil {
		return fmt.Errorf("failed to Visit http://openinsider.com/latest-insider-sales-100k")
	}

	return nil
}

func isToday(filingDate string) bool {
	parsedDate := strings.Split(strings.Split(filingDate, " ")[0], "-")
	date := Map(parsedDate, strconv.Atoi)

	now := time.Now()
	if now.Day()-1 == date[2] {
		return true
	}
	return false
}

func Map(slice []string, mapper func(string) (int, error)) []int {
	mappedSlice := make([]int, len(slice))
	for i, v := range slice {
		result, err := mapper(v)
		if err != nil {
			panic(err)
		}
		mappedSlice[i] = result
	}
	return mappedSlice
}
