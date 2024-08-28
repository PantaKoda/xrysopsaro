package common

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func RequestBody(urlString string) (*goquery.Document, error) {

	res, err := http.Get(urlString)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}
