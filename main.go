package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func main() {
	// Create a new Colly collector
	c := colly.NewCollector()

	// Define an endpoint to scrape data
	// curl http://localhost:8080/scrape
	http.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {

		targetURL := "https://www.infobae.com/"

		var scrapedData []string

		// Set up the callback for when a scraped element is found
		c.OnHTML("div.top-home", func(e *colly.HTMLElement) {
			scrapedData = append(scrapedData, e.Text)
		})

		// Start the scraping process
		if err := c.Visit(targetURL); err != nil {
			log.Println("Error visiting the website:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(scrapedData)
		if err != nil {
			log.Println("Error encoding JSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	// Start the HTTP server
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
