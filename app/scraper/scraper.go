package scraper

import (
	"backend/app/helpers"
	"backend/app/model"
	"fmt"
	"log"
	"time"

	// "time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Run() (*model.Entry, error) {
	var obj *model.Entry

	// Step 1: Cek apakah sudah ada data pada `date` tertentu
	// Step 2: Jika belum, mulai web scraping
	// Step 3: Lakukan validasi
	// Step 4: Insert to DB
	// Step 5: Response 201

	// Instantiate default collector as "c"
	c := colly.NewCollector(
		// Domains to be visited: www.bca.co.id
		colly.AllowedDomains("www.bca.co.id"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		// colly.CacheDir("./cache"),
	)

	// Run callback upon finding the exchange rate table element
	c.OnHTML("table.m-table-kurs", func(e *colly.HTMLElement) {
		var arr []model.Exc_rates
		// ForEach will execute callback on every matching element
		// inside element "e", collecting data for each row
		e.ForEach("tbody tr", func(i int, row *colly.HTMLElement) {

			// Scraping results in string data. Converting string number to float32 will
			// force the code to estimate value for the digits behind decimal
			rowData := model.Exc_rates{
				Symbol: row.Attr("code"),
				Rates: model.Rates{
					ER: model.Buysell{
						Buy:  row.ChildText("p[rate-type='ERate-buy']"),
						Sell: row.ChildText("p[rate-type='ERate-sell']"),
					},
					TT: model.Buysell{
						Buy:  row.ChildText("p[rate-type='TT-buy']"),
						Sell: row.ChildText("p[rate-type='TT-sell']"),
					},
					BN: model.Buysell{
						Buy:  row.ChildText("p[rate-type='BN-buy']"),
						Sell: row.ChildText("p[rate-type='BN-sell']"),
					},
				},
			}
			arr = append(arr, rowData)
		})
		currentTime := time.Now()
		stringTime := currentTime.Format("2006-01-02")
		dateTime, err := helpers.YYYYMMDDToISODate(stringTime)

		if err != nil {
			// Throw error
		}

		obj = &model.Entry{
			Id:         primitive.NewObjectID(),
			Date:       dateTime,
			Exc_Rates:  arr,
			Created_At: time.Now(),
			Updated_At: time.Now(),
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	fmt.Println("Scraping page...")
	c.Visit("https://www.bca.co.id/en/informasi/kurs")

	log.Printf("Scraping Complete!\n")
	log.Println(c)

	return obj, nil
}
