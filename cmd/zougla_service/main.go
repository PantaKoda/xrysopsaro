package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func PostId(gui_url string) string {

	u, err := url.Parse(gui_url)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	return q.Get("p")

}

func GetImageUrl(post_link string) string {

	doc, err := RequestBody(post_link)

	if err != nil {
		log.Fatal("Couldn't parse body doc")
	}

	imgUrl, exists := doc.Find("main .left-part .entry-img img").Attr("src")

	if !exists {
		return ""
	}

	return imgUrl
}

func ZouglaTimeToRFC1123Z(pub_date string) string {

	// Parse the input string
	t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", pub_date)
	if err != nil {
		return ""
	}

	// Load the Athens timezone
	athensLocation, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		return ""
	}

	// Create a new time with the same UTC time, but in Athens timezone
	athensTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), athensLocation)

	// Format the time in RFC3339 format
	return athensTime.Format(time.RFC3339)

}

func PrintPost(RssFeed RssZougla) {

	for _, item := range RssFeed.Channel.Item {

		ImgUrl := GetImageUrl(item.Link)

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("PubDateFormatted:", ZouglaTimeToRFC1123Z(item.PubDate))
		fmt.Println("Categories:", strings.Join(item.Category, ","))
		//fmt.Println("Guid:", item.Guid.Text)
		fmt.Println("PostId :", PostId(item.Guid.Text))
		fmt.Println("ImgUrl  :", ImgUrl)

		fmt.Println()
	}
}

func main() {

	var RssFeed RssZougla

	err := ParseFeed("https://zougla.gr/feed/", &RssFeed)

	if err != nil {
		fmt.Println("Error in XML PARSE")
		return
	}

	PrintPost(RssFeed)
}
