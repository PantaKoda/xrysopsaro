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

func ParseTime(timeStr string) time.Time {
	// Step 1: Parse the RFC1123Z formatted time string.
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Fatalf("Error parsing time: %v", err)
	}
	return parsedTime
}

func PrintPost(RssFeed RssNewsBeast) {

	for _, item := range RssFeed.Channel.Item {

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", strings.TrimSpace(item.Description))
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("Categories:", strings.Join(item.Category, ","))
		fmt.Println("Image:", item.Enclosure.URL)

		fmt.Println()
	}
}

func main() {

	var RssFeed RssNewsBeast

	err := ParseFeed("https://www.newsbeast.gr/feed", &RssFeed)

	if err != nil {
		fmt.Println("Error in XML PARSE", err)
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
				ImgUrl:         pgtype.Text{String: item.Thumbnail.URL, Valid: true},
				Categories:     pgtype.Text{String: strings.Join(item.Category, ","), Valid: true},
				Url:            item.Link,
				Website:        "newsbeast",
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
