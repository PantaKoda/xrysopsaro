package main

import (
	"context"
	"fmt"
	cache2 "github.com/PantaKoda/xrysopsaro/common/cache"
	database2 "github.com/PantaKoda/xrysopsaro/common/database_sqlc"
	"github.com/PantaKoda/xrysopsaro/common/dbconnect"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const DEFAULT_FILENAME = "LOCAL_CACHE.txt"

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

// Converts an RFC1123Z UTC time string to RFC3339 format with Athens local time zone.
func ParseTime(timeStr string) time.Time {
	// Step 1: Parse the RFC1123Z formatted time string.
	parsedTime, err := time.Parse(time.RFC1123Z, timeStr)
	if err != nil {
		log.Fatalf("Error parsing time: %v", err)
	}

	// Step 2: Load the Athens time zone.
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		log.Fatalf("Error loading time zone: %v", err)
	}

	// Step 3: Convert the parsed time to Athens time zone.
	localTime := parsedTime.In(location)

	// Step 4: Format the time to RFC3339, ensuring the time zone is correct.
	return localTime
}

func PrintPost(RssFeed RssZougla) {

	for _, item := range RssFeed.Channel.Item {

		ImgUrl := GetImageUrl(item.Link)

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("PubDateFormatted:", ParseTime(item.PubDate))
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

	// Initialize database_sqlc connection
	conn := dbconnect.InitializeDatabase()
	defer conn.Close(context.Background())

	// Create an instance of Queries using the connection
	queries := database2.New(conn)

	cachedUrls, err := cache2.ReadFromCache(DEFAULT_FILENAME)

	if err != nil {
		log.Fatal(err)
	}

	articlesUrlsList := []string{}

	for _, item := range RssFeed.Channel.Item {

		if !slices.Contains(cachedUrls, item.Link) {

			ImgUrl := GetImageUrl(item.Link)

			articlesUrlsList = append(articlesUrlsList, item.Link)

			// CreatePostParams is from SQLC, mapped with the scraped data
			post := database2.CreatePostParams{

				Title:          item.Title,
				PublishDate:    pgtype.Timestamptz{Time: ParseTime(item.PubDate), Valid: true},
				PublishDateRaw: item.PubDate,
				Description:    pgtype.Text{String: item.Description, Valid: true},
				ImgUrl:         pgtype.Text{String: ImgUrl, Valid: true},
				Categories:     pgtype.Text{String: strings.Join(item.Category, ","), Valid: true},
				Url:            item.Link,
				Website:        "zougla",
			}
			// Insert into the database_sqlc using SQLC
			_, err := queries.CreatePost(context.Background(), post)
			if err != nil {
				log.Printf("Failed to insert post into the database_sqlc: %v", err)
			} else {
				log.Printf("Inserted post: %s", item.Title)
			}
		}

	}

	fmt.Println("Total URLs to append:", len(articlesUrlsList)) // Debug statement§

	err = cache2.AppendSliceToFile(DEFAULT_FILENAME, articlesUrlsList)

	if err != nil {
		log.Fatal("Error while appending to file", err)
	}

	//PrintPost(RssFeed)
}
