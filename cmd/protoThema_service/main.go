package main

//go build -o ./output_test -a -ldflags="-s -w"  main.go ReadFromCache.go CurrentDirectory.go AppendToCache.go
import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/PantaKoda/xrysopsaro/pkg/database"
	"github.com/PuerkitoBio/goquery"
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

func ProtoThemaTimeToRFC3339(pub_date string) string {

	// Parse the time in RFC 3339 format
	parsedTime, err := time.Parse(time.RFC3339, pub_date)

	if err != nil {
		fmt.Println("Error parsing time:", err)
		return ""
	}

	return parsedTime.Format(time.RFC3339)

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

type Post struct {
	Title       string `json:"title"`
	PublishDate string `json:"publish_date"`
	Description string `json:"description"`
	ImgUrl      string `json:"img_url"`
	Categories  string `json:"categories"`
	Url         string `json:"url"`
	Website     string `json:"website"`
}

func main() {

	// Initialize the database connection
	dbConn := database.InitializeDatabase()

	// Set the connection as the singleton instance
	database.InitializeDB(dbConn)

	// Ensure the connection is closed at the end of the program
	defer func() {
		if err := database.GetDB().Close(); err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}()

	// Example: Just to ensure the connection is working
	if err := database.PingDatabase(database.GetDB()); err != nil {
		log.Fatal(err)
	}
	log.Println("Pinged database successfully")

	doc, err := RequestBody(BASE_URL)

	if err != nil {

		log.Fatal("Error fetching page source")
	}

	articlesListSelector := ARTICLES_LIST_SELECTOR
	sel := doc.Find(articlesListSelector)
	cachedUrls, err := ReadFromCache(DEFAULT_FILENAME)

	if err != nil {
		log.Fatal(err)
	}

	articlesUrlsList := []string{}
	articlesList := []Post{}

	for i := range sel.Nodes {

		single := sel.Eq(i)

		postUrl, exists := GetPostUrl(single)

		if !exists {
			continue
		}

		if IsThrowablePort(single) {
			continue
		}

		var post Post

		if !slices.Contains(cachedUrls, postUrl) {

			articlesUrlsList = append(articlesUrlsList, postUrl)

			fmt.Println("Appending URL to list:", postUrl) // Debug statement
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

			post = Post{
				Title:       title,
				PublishDate: publishedDate,
				Description: description,
				ImgUrl:      imgUrl,
				Categories:  tags,
				Url:         postUrl,
				Website:     "protoThema",
			}
			articlesList = append(articlesList, post)
		}

	}

	//TODO Save to Database

	PrintPostList(articlesList)

	fmt.Println("Total URLs to append:", len(articlesUrlsList)) // Debug statement

	err = appendSliceToFile(DEFAULT_FILENAME, articlesUrlsList)

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

}
