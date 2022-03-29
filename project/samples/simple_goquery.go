package samples

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func TestScrape() {
	// Request the html page.
	res, err := http.Get("https://finance.yahoo.com/quote/ORCL")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load the html document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find the review items
	doc.Find("#quote-header-info").Each(func(i int, s *goquery.Selection) {
		// for each item found, get the name.
		name := s.Find("h1").Text()
		fmt.Println(name)
	})
}
