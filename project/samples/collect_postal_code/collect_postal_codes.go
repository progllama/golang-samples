package collectpostalcode

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func CollectCountryNamesByGoquery() {
	url := "https://www.mofa.go.jp/mofaj/area/index.html"
	countryNames := collectCountryNames(url)
	save(countryNames)
}

func collectCountryNames(url string) []string {
	doc, err := createtDocument(url)
	if err != nil {
		log.Fatal(err)
		return []string{}
	}

	return collectCountryNamesFromDoc(doc)
}

func createtDocument(url string) (*goquery.Document, error) {
	if hasCache(url) {
		log.Println("From cache.")
		return createtDocumentFromCache(url)
	} else {
		log.Println("From web.")
		return createtDocumentFromWebResource(url)
	}
}

func collectCountryNamesFromDoc(doc *goquery.Document) []string {
	countryNames := []string{}
	doc.Find("li.styled2").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, t *goquery.Selection) {
			countryNames = append(countryNames, t.Text())
		})
	})
	return countryNames
}

func save(lines []string) error {
	f, err := os.Create(getSavePath())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	err = writeLines(f, lines)
	if err != nil {
		return err
	}

	return nil
}

func writeLines(f *os.File, lines []string) error {
	for _, line := range lines {
		err := writeLine(f, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLine(f *os.File, line string) error {
	_, err := f.WriteString(line + "\n")
	return err
}

func getSavePath() string {
	return "./CountryNames.txt"
}

func createtDocumentFromCache(url string) (*goquery.Document, error) {
	r, err := getCacheReader(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(r)
}

func createtDocumentFromWebResource(url string) (*goquery.Document, error) {
	resp, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	doc, err := createDocumentFromResponse(url, resp)
	if err != nil {
		return &goquery.Document{}, err
	}
	return doc, nil
}

func createDocumentFromResponse(url string, resp *http.Response) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return doc, err
	}
	createCacheFromDoc(url, doc)
	return doc, nil
}

func createCacheFromDoc(url string, doc *goquery.Document) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(getCachePath(url), []byte(html), os.ModePerm)
	if err != nil {
		return err
	}
	return err
}

func sendRequest(url string) (*http.Response, error) {
	return http.Get(url)
}

func getCacheReader(url string) (io.Reader, error) {
	b, err := ioutil.ReadFile(getCachePath(url))
	return bytes.NewBuffer(b), err
}

func hasCache(url string) bool {
	dest := getCachePath(url)
	return isPathExist(dest)
}

func isPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getCachePath(url string) string {
	base := "./"
	hash := getHashString(url)
	return base + hash + ".html"
}

func getHashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
