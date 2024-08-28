package main

import (
	"context"
	"fmt"
	"github.com/PantaKoda/xrysopsaro/common"
	"github.com/PantaKoda/xrysopsaro/common/cache"
	"github.com/PantaKoda/xrysopsaro/common/database_sqlc"
	"github.com/PantaKoda/xrysopsaro/common/dbconnect"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/url"
	"slices"
	"strings"
	"time"
)

const (
	URL              = "https://www.news247.gr/roi-eidiseon/"
	SELECTOR         = "#latest_news_header .latest_news_article_container article"
	DEFAULT_FILENAME = "LOCAL_CACHE.txt"
)

func IsThrowablePort(node *goquery.Selection) bool {

	ExternalSite, exists := node.Find(".post__content h3.post__title a").Attr("href")

	u, err := url.Parse(ExternalSite)

	if err != nil {
		log.Fatal("Couldnt parse external site href", err)
	}

	hostname := u.Hostname()

	if !exists || !strings.Contains(hostname, "news247.gr") {
		return true
	}

	return false
}

func main() {
	doc, err := common.RequestBody(URL)

	if err != nil {
		log.Fatal("Cannot fetch source html")
	}

	sel := doc.Find(SELECTOR)
	cachedUrls, err := cache.ReadFromCache(DEFAULT_FILENAME)
	if err != nil {
		log.Fatal("Failed reading cache", err)
	}
	articlesUrlsList := []string{}

	// Initialize database_sqlc connection
	conn := dbconnect.InitializeDatabase()
	queries := database_sqlc.New(conn)

	for i := range sel.Nodes {

		single := sel.Eq(i)
		//log.Println(sel.Eq(i).Html())

		if IsThrowablePort(single) {
			continue
		}

		postUrl, exists := GetPostUrl(single)

		if !exists {
			log.Fatal("Couldn't find post url")
		}

		if !slices.Contains(cachedUrls, postUrl) {

			articleBody, err := common.RequestBody(postUrl)

			if err != nil {
				log.Fatal("Cannot fetch article's source html")
			}

			articlesUrlsList = append(articlesUrlsList, postUrl)

			title := GetTitle(single)
			description := GetDescription(articleBody)
			tags := GetArticlesTags(articleBody)
			formattedDate, publishDateRaw := GetPublishDate(single)
			imgUrlBig := GetImageUrl(articleBody)

			// CreatePostParams is from SQLC, mapped with the scraped data
			post := database_sqlc.CreatePostParams{

				Title:          title,
				PublishDate:    pgtype.Timestamptz{Time: formattedDate, Valid: true},
				PublishDateRaw: publishDateRaw,
				Description:    pgtype.Text{String: description, Valid: true},
				ImgUrl:         pgtype.Text{String: imgUrlBig, Valid: true},
				Categories:     pgtype.Text{String: tags, Valid: true},
				Url:            postUrl,
				Website:        "news247",
			}

			// Insert into the database_sqlc using SQLC
			_, err = queries.CreatePost(context.Background(), post)

			if err != nil {
				log.Printf("Failed to insert post into the database_sqlc: %v", err)
			} else {
				log.Printf("Inserted post: %s", title)
			}

		}

	}

	fmt.Println("Total URLs to append:", len(articlesUrlsList)) // Debug statement§

	err = cache.AppendSliceToFile(DEFAULT_FILENAME, articlesUrlsList)

	if err != nil {
		log.Fatal("Error while appending to file", err)
	}

}

func GetImageUrl(node *goquery.Document) string {
	//small single *goquery.Selection
	//k, exists := single.Find("figure a img").Attr("data-src")
	l, exists := node.Find(".single_article__header .single_article__main_image img").Attr("src")

	if !exists {
		return ""
	}

	return l
}

func GetPublishDate(node *goquery.Selection) (time.Time, string) {

	//28.08.2024 22:55 Or 00:00
	t := strings.TrimSpace(node.Find(".post__content .caption.post__date.article-xs-font").Text())
	//k := articleBody.Find(".single_article__header_info .article__date.s-font span").Text()

	// Load the Athens time zone (Europe/Athens)
	location, err := time.LoadLocation("Europe/Athens")

	if err != nil {
		fmt.Println("Failed to load Athens timezone:", err)
		return time.Time{}, t
	}

	// Check if it's only the time, e.g., "00:00"
	if len(t) == 5 && strings.Contains(t, ":") {
		// Add today's date to the time
		currentDate := time.Now().In(location).Format("02.01.2006")
		t = currentDate + " " + t // e.g., "28.08.2024 00:00"
	}

	parseTime, err := time.ParseInLocation("02.01.2006 15:04", t, location)
	if err != nil {
		log.Fatal("Failed to parse date:", err)
		return time.Time{}, t
	}

	return parseTime, t
}

func GetArticlesTags(articleBody *goquery.Document) string {
	tagsListHtml := articleBody.Find(".article__tags_container .article__tags ul .tag_item")

	tagsList := []string{}

	tagsListHtml.Each(func(i int, item *goquery.Selection) {

		tagText := strings.TrimSpace(item.Text())
		tagsList = append(tagsList, tagText)
	})
	return strings.Join(tagsList, ",")
}

func GetDescription(articleBody *goquery.Document) string {

	description := articleBody.Find(".single_article__main .article__lead p").Text()

	return strings.TrimSpace(description)
}

func GetTitle(node *goquery.Selection) string {
	return node.Find(".post__content h3.post__title a").Text()
}

func GetPostUrl(node *goquery.Selection) (string, bool) {

	articleUrl, exists := node.Find(".post__content h3.post__title a").Attr("href")

	return articleUrl, exists
}

/*func PrintPost() {
	fmt.Println("Title : ", title)
	fmt.Println("description : ", description)
	fmt.Println("postUrl : ", postUrl)
	fmt.Println("tags : ", tags)
	fmt.Println("Original date : ", publishDateRaw)
	fmt.Println("Formatted date : ", formattedDate.Format(time.RFC3339))

	fmt.Println("Image Url big : ", imgUrlBig)
	fmt.Println()
}*/
