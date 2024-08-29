package main

import (
	"context"
	"fmt"
	cache2 "github.com/PantaKoda/xrysopsaro/common/cache"
	database2 "github.com/PantaKoda/xrysopsaro/common/database_sqlc"
	"github.com/PantaKoda/xrysopsaro/common/dbconnect"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"slices"
	"strings"
	"time"
)

const DEFAULT_FILENAME = "LOCAL_CACHE.txt"

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

func main() {
	var RssFeed RssEnikos

	err := ParseFeed("https://enikos.gr/feed/", &RssFeed)
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

			post := database2.CreatePostParams{

				Title:          item.Title,
				PublishDate:    pgtype.Timestamptz{Time: ParseTime(item.PubDate), Valid: true},
				PublishDateRaw: item.PubDate,
				Description:    pgtype.Text{String: item.Description, Valid: true},
				ImgUrl:         pgtype.Text{String: item.Image, Valid: true},
				Categories:     pgtype.Text{String: strings.Join(item.Category, ","), Valid: true},
				Url:            item.Link,
				Website:        "enikos",
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

func PrintPost(RssFeed RssEnikos) {

	for _, item := range RssFeed.Channel.Item {

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Creator:", item.Creator)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("PubDate Formatted 2 :", ParseTime(item.PubDate).Format(time.RFC3339))
		fmt.Println("Categories:", strings.Join(item.Category, ","))
		fmt.Println("Image:", item.Image)

		fmt.Println()
	}
}
