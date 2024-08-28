package main

import (
	"fmt"
	"strings"
)

func PrintPost(RssFeed RssEnikos) {

	for _, item := range RssFeed.Channel.Item {

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Creator:", item.Creator)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("Categories:", strings.Join(item.Category, ","))
		fmt.Println("Image:", item.Image)

		fmt.Println()
	}
}

func main() {
	var RssFeed RssEnikos

	err := ParseFeed("https://enikos.gr/feed/", &RssFeed)
	if err != nil {
		fmt.Println("Error in XML PARSE")
		return
	}

	PrintPost(RssFeed)

}
