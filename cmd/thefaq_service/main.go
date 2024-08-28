package main

import (
	"context"
	"fmt"
	cache2 "github.com/PantaKoda/xrysopsaro/common/cache"
	database2 "github.com/PantaKoda/xrysopsaro/common/database_sqlc"
	"github.com/PantaKoda/xrysopsaro/common/dbconnect"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
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

func GetImgUrl(articleUrl string, articleId string) string {
	doc, err := RequestBody(articleUrl)
	if err != nil {
		log.Fatal("Error fetching news feed source:", err)
	}

	// Debugging: print the HTML to understand its structure.
	// fmt.Println(doc.Html())

	// Find the first image tag with class 'wp-post-image' (modify selector as needed).
	selector := "#post-" + articleId + " > img"
	sel := doc.Find(selector)
	//#post-306522 > img

	// Try to get the image URL from various possible attributes.
	imgURL, exists := sel.Attr("data-breeze")
	if !exists {
		imgURL, exists = sel.Attr("data-brsrcset")
		if exists {
			// Extract the first URL from srcset by splitting the string.
			imgURL = strings.Split(imgURL, " ")[0]
		} else {
			// Fall back to the src attribute.
			imgURL, exists = sel.Attr("src")
		}
	}

	// Check if an image URL was found; if not, set an empty string.
	if !exists {
		imgURL = ""
	}

	return imgURL
}

// Function to extract text without HTML tags
func ExtractTextFromHTML(htmlString string) string {
	// Create a new goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if err != nil {
		return ""
	}

	// Extract the plain text from the document
	text := doc.Text()

	// Clean up any excessive whitespace
	text = strings.TrimSpace(text)

	return text
}

func main() {

	var RssFeed RssFaq

	err := ParseFeed("https://thefaq.gr/feed/", &RssFeed)
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

			articlesUrlsList = append(articlesUrlsList, item.Link)
			u, error := url.Parse(item.Guid.Text)
			if error != nil {
				log.Fatal("Cannot parse article gui id", err)
			}

			q := u.Query()
			articleId := q.Get("p")

			post := database2.CreatePostParams{

				Title:          strings.TrimSpace(item.Title),
				PublishDate:    pgtype.Timestamptz{Time: ParseTime(item.PubDate), Valid: true},
				PublishDateRaw: item.PubDate,
				Description:    pgtype.Text{String: strings.TrimSpace(ExtractTextFromHTML(item.Description)), Valid: true},
				ImgUrl:         pgtype.Text{String: GetImgUrl(item.Link, articleId), Valid: true},
				Categories:     pgtype.Text{String: strings.Join(item.Category, ","), Valid: true},
				Url:            item.Link,
				Website:        "thefaq",
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

}
