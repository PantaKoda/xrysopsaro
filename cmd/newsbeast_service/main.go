package main

import (
	"fmt"
	"regexp"
	"strings"
)

func PostId(post_link string) string {

	re := regexp.MustCompile(`/(\d+)/`)

	// Find the first match in the URL
	match := re.FindStringSubmatch(post_link)

	// Check if we found a match
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func PrintPost(RssFeed RssNewsBeast) {

	for _, item := range RssFeed.Channel.Item {

		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", strings.TrimSpace(item.Description))
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println("Categories:", strings.Join(item.Category, ","))
		fmt.Println("Image:", item.Enclosure.URL)
		fmt.Println("PostID : ", PostId(item.Link))

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

	PrintPost(RssFeed)

}
