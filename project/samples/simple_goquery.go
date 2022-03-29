package samples

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
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

type Ticker struct {
	Name   string
	Symbol string
	Price  float64
	URL    string
}

func readCSVFile(filePath string) map[string]Ticker {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file"+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	tickers := make(map[string]Ticker)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		url := "https://finance.yahoo.com/quote/" + record[0]
		var ticker = Ticker{
			Name:   record[1],
			Symbol: record[0],
			URL:    url,
		}
		tickers[record[0]] = ticker
	}

	if err != nil {
		log.Fatal("unable to parse file as csv for "+filePath, err)
	}

	return tickers
}

func TestReadCSVFile() {
	start := time.Now()
	tickers := readCSVFile("./test.csv")

	c := colly.NewCollector()

	c.OnHTML("#quote-header-info", func(e *colly.HTMLElement) {
		name := e.ChildText("h1")
		quote := e.ChildTexts("span")

		temp := strings.Split(name, "(")
		name = temp[0]
		symbol := temp[1][:len(temp[1])-1]

		price, _ := strconv.ParseFloat(quote[3], 32)

		ticker := tickers[symbol]
		ticker.Name = name
		ticker.Price = price

		tickers[symbol] = ticker
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	count := 0
	for _, ticker := range tickers {
		count++
		c.Visit(ticker.URL)
		if count == 15 {
			break
		}
	}

	fmt.Println(time.Since(start))
}
