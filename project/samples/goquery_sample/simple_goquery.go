package goquery_sample

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

// Using Version
// wget https://chromedriver.storage.googleapis.com/100.0.4896.60/chromedriver_linux64.zip
func setupDriver(url string) string {
	var driver *agouti.WebDriver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		log.Fatal(err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// URL Setting
	err = page.Navigate(url)
	if err != nil {
		log.Fatal(err)
	}

	src, err := page.HTML()
	if err != nil {
		log.Fatal(err)
	}
	return src
}

type WebSpider struct {
	url      string
	keyword  string
	document *goquery.Document
}

func NewWebSpider(url string, keyword string) *WebSpider {
	return &WebSpider{
		url:      url,
		keyword:  keyword,
		document: nil,
	}
}

func (ws *WebSpider) Run() {
	ws.document, _ = ws.buildDocument()
	route := ws.findRoutesFromKeyWord()[ws.keyword]
	ws.document.Find(strings.Join(route, " ")).Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}

func (ws *WebSpider) buildDocument() (*goquery.Document, error) {
	// response, err := http.Get(ws.url)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return &goquery.Document{}, err
	// }

	srcPage := setupDriver(ws.url)
	sr := strings.NewReader(srcPage)

	document, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		log.Fatal(err)
		return &goquery.Document{}, err
	}

	return document, nil
}

// Routesと複数形の理由はキーワードに正規表現を使ったときを想定して。
func (ws *WebSpider) findRoutesFromKeyWord() map[string][]string {
	routes := make(map[string][]string, 0)
	body := ws.document.Find("body")
	body.Find("*").Each(func(i int, s *goquery.Selection) {
		text := s.Clone().Children().Remove().End().Text()
		if text == ws.keyword {
			routes[text] = ws.findRouteFromNode(s)
		}
	})
	return routes
}

func (ws *WebSpider) findRouteFromNode(node *goquery.Selection) []string {
	route := make([]string, 0)
	currentNode := node
	for {
		route = append(route, goquery.NodeName(currentNode))
		if goquery.NodeName(currentNode) == "html" {
			break
		}
		currentNode = currentNode.Parent()
	}
	// reverse
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}
	return route
}

func GoquerySample() {
	ws := NewWebSpider(
		"https://ja.wikipedia.org/wiki/%E5%9B%BD%E3%81%AE%E4%B8%80%E8%A6%A7_(%E5%A4%A7%E9%99%B8%E5%88%A5)",
		"バクー",
	)
	ws.Run()
}

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

// func TestReadCSVFile() {
// 	start := time.Now()
// 	tickers := readCSVFile("./test.csv")

// 	c := colly.NewCollector()

// 	c.OnHTML("#quote-header-info", func(e *colly.HTMLElement) {
// 		name := e.ChildText("h1")
// 		quote := e.ChildTexts("span")

// 		temp := strings.Split(name, "(")
// 		name = temp[0]
// 		symbol := temp[1][:len(temp[1])-1]

// 		price, _ := strconv.ParseFloat(quote[3], 32)

// 		ticker := tickers[symbol]
// 		ticker.Name = name
// 		ticker.Price = price

// 		tickers[symbol] = ticker
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})

// 	count := 0
// 	for _, ticker := range tickers {
// 		count++
// 		c.Visit(ticker.URL)
// 		if count == 15 {
// 			break
// 		}
// 	}

// 	fmt.Println(time.Since(start))
// }

var tickerSymbols = make([]string, 0)

var maxGoroutines int = 10
var tickersToBeScraped int = 2000

// func TestReadGoScraper() {
// 	start := time.Now()
// 	tickers := readCSVFile("../hoge.csv")

// 	c := colly.NewCollector()

// 	var wg sync.WaitGroup

// 	guard := make(chan struct{}, maxGoroutines)

// 	c.OnHTML("#quote-header-info", func(e *colly.HTMLElement) {
// 		name := e.ChildText("h1")
// 		quote := e.ChildTexts("span")

// 		temp := strings.Split(name, "(")
// 		name = temp[0]
// 		symbol := temp[1][:len(temp[1])-1]

// 		price, _ := strconv.ParseFloat(quote[3], 32)
// 		price = math.Round(price/0.01) * 0.01

// 		ticker := tickers[symbol]
// 		ticker.Name = name
// 		ticker.Price = price

// 		tickers[symbol] = ticker
// 		fmt.Println(price)
// 	})

// 	c.OnError(func(_ *colly.Response, err error) {
// 		log.Println("Something went wrong.", err)
// 		<-guard
// 		wg.Done()
// 	})

// 	c.OnResponse(func(r *colly.Response) {
// 		<-guard
// 		wg.Done()
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})

// 	sort.Strings(tickerSymbols)
// 	for i := 0; i < tickersToBeScraped; i++ {
// 		guard <- struct{}{}
// 		wg.Add(1)
// 		go c.Visit(tickers[tickerSymbols[i]].URL)
// 	}

// 	wg.Wait()
// 	count := 0
// 	for i := 0; i < tickersToBeScraped; i++ {
// 		if tickers[tickerSymbols[i]].Price != 0.00 {
// 			count++
// 		}
// 	}
// 	fmt.Println("Data Successfully scraped for:")
// 	fmt.Println(count, "/", tickersToBeScraped)
// 	fmt.Println(time.Since(start))
// }
