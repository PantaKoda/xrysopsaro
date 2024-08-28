package main

//go build -o ./output_test -a -ldflags="-s -w"  main.go ReadFromCache.go CurrentDirectory.go AppendToCache.go
import (
	"context"
	"fmt"
	cache2 "github.com/PantaKoda/xrysopsaro/common/cache"
	database2 "github.com/PantaKoda/xrysopsaro/common/database_sqlc"
	"github.com/PantaKoda/xrysopsaro/common/dbconnect"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	DEFAULT_FILENAME       = "LOCAL_CACHE.txt"
	ARTICLES_LIST_SELECTOR = ".outer .inner .mainWrp .main .content .listLoop article"
	BASE_URL               = "https://www.protothema.gr/oles-oi-eidiseis/"
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

// Modified parseTime function to ensure the original time zone is preserved.
func parseTime(timeStr string) time.Time {
	// Parse the time string using RFC3339 to maintain the original timezone information.
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return time.Now() // Default fallback, but ideally handle errors properly.
	}
	return parsedTime
}

func TagsCategoriesProtoThema(PostUrl string) string {

	doc, err := RequestBody(PostUrl)

	if err != nil {

		log.Fatal("Error fetching page source")
	}

	selector := ".outer main > section.mainSection .articleTopInfo .tagsCnt a"

	/* if strings.Contains(PostUrl, "sports") {
		selector = ".inner .articleTopInfo .tagsWrp a"
	} */

	tagsList := []string{}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {

		tag := strings.TrimSpace(s.Text())
		tagsList = append(tagsList, tag)

	})

	return strings.Join(tagsList, ",")
}

func IsThrowablePort(node *goquery.Selection) bool {

	isEnglishNews := node.Find(".desc > a").First().Text()

	if isEnglishNews == "English News" {
		return true
	}

	ExternalSite, exists := node.Find("a").First().Attr("href")

	u, err := url.Parse(ExternalSite)

	if err != nil {
		log.Fatal("Couldnt parse external site href", err)
	}

	hostname := u.Hostname()

	if !exists || !strings.Contains(hostname, "protothema.gr") {
		return true
	}

	return false
}

func main() {

	// Initialize database_sqlc connection
	conn := dbconnect.InitializeDatabase()
	defer conn.Close(context.Background())

	// Create an instance of Queries using the connection
	queries := database2.New(conn)

	doc, err := RequestBody(BASE_URL)

	if err != nil {

		log.Fatal("Error fetching page source")
	}

	sel := doc.Find(ARTICLES_LIST_SELECTOR)
	cachedUrls, err := cache2.ReadFromCache(DEFAULT_FILENAME)

	if err != nil {
		log.Fatal(err)
	}

	articlesUrlsList := []string{}

	for i := range sel.Nodes {

		single := sel.Eq(i)

		postUrl, exists := GetPostUrl(single)

		if !exists {
			continue
		}

		if IsThrowablePort(single) {
			continue
		}

		if !slices.Contains(cachedUrls, postUrl) {

			articlesUrlsList = append(articlesUrlsList, postUrl)

			//fmt.Println("Appending URL to list:", postUrl) // Debug statement
			title := GetTitle(single)
			description := GetDescription(single)

			publishedDate, exists := GetPublishDate(single)

			if !exists {
				publishedDate = ""
			}

			imgUrl, exists := getImgUrl(single)
			if !exists {
				imgUrl = ""
				log.Println("Cannot find Image Url")
			}

			tags := TagsCategoriesProtoThema(postUrl)

			// CreatePostParams is from SQLC, mapped with the scraped data
			post := database2.CreatePostParams{

				Title:          title,
				PublishDate:    pgtype.Timestamptz{Time: parseTime(publishedDate), Valid: true},
				PublishDateRaw: publishedDate,
				Description:    pgtype.Text{String: description, Valid: true},
				ImgUrl:         pgtype.Text{String: imgUrl, Valid: true},
				Categories:     pgtype.Text{String: tags, Valid: true},
				Url:            postUrl,
				Website:        "protoThema",
			}
			// Insert into the database_sqlc using SQLC
			_, err := queries.CreatePost(context.Background(), post)
			if err != nil {
				log.Printf("Failed to insert post into the database_sqlc: %v", err)
			} else {
				log.Printf("Inserted post: %s", title)
			}

			//articlesList = append(articlesList, post)
		}

	}

	//PrintPostList(articlesList)

	fmt.Println("Total URLs to append:", len(articlesUrlsList)) // Debug statement§

	err = cache2.AppendSliceToFile(DEFAULT_FILENAME, articlesUrlsList)

	if err != nil {
		log.Fatal("Error while appending to file", err)
	}
}

func GetTitle(node *goquery.Selection) string {

	title := node.Find(".desc .heading h3 a").Text()
	return strings.TrimSpace(title)
}

func GetDescription(node *goquery.Selection) string {

	description := node.Find(".desc .txt").Text()
	return strings.TrimSpace(description)
}

func GetPostUrl(node *goquery.Selection) (string, bool) {

	return node.Find("a").Attr("href")
}

func GetPublishDate(node *goquery.Selection) (string, bool) {
	return node.Find(".wrp time").Attr("datetime")
}

func getImgUrl(node *goquery.Selection) (string, bool) {

	return node.Find("figure a picture img").Attr("data-src")
}

/*
	html, error := fmt.Println(node.Html())
	if error != nil {
	}
	fmt.Println("INSIDE GETIMG URP", html)
*/
/*
func PrintPostList(articleList []Post) {
	for i := range articleList {
		articleList[i].PrintDetails()
	}
}

func (p *Post) PrintDetails() {

	fmt.Println("Title :", p.Title)
	fmt.Println("Description :", p.Description)
	fmt.Println("Publish Date :", p.PublishDate)
	fmt.Println("Image URL :", p.ImgUrl)
	fmt.Println("Categories :", p.Categories)
	fmt.Println("Url :", p.Url)
	fmt.Print()

} */
