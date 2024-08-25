package main

import (
	"fmt"
)

func main() {
	var RssFeed Rss

	err := ParseFeed("https://enikos.ge/feed/", &RssFeed)
	if err != nil {
		fmt.Println("Error in XML PARSE")
		return
	}

	for _, item := range RssFeed.Channel.Item {
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Creator:", item.Creator)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("Categories:", item.Category)
		fmt.Println("Image:", item.Image)

		fmt.Println()
	}
}
