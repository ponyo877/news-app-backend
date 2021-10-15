package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/mmcdole/gofeed"
	"github.com/ponyo877/news-app-backend/handler/elastic"
	"github.com/ponyo877/news-app-backend/handler/imagectl"
)

// SiteInfo is metainfomation of RSS Site
type SiteInfo struct {
	ID         int
	title      string
	rssURL     string
	latestDate string
}

// SiteRecord is article infomation for DB
type SiteRecord struct {
	title      string
	URL        string
	image      string
	updateDate string
	siteID     int
}

func PutPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		update_count := PutPostTmp()
		return c.JSON(http.StatusOK, map[string]int{"update_count": update_count})
	}
}

func PutPostJob() {
	update_count := PutPostTmp()
	fmt.Println("update_count:", update_count)
}

func PutPostTmp() int {
	siteinfolist := getSiteInfoList()
	feedparser := gofeed.NewParser()
	feedArray := []SiteRecord{}
	isVisit := map[int]bool{}
	update_count := 0
	for _, siteinfo := range siteinfolist {
		isVisit[siteinfo.ID] = false
		feed, _ := feedparser.ParseURL(siteinfo.rssURL)
		items := feed.Items
		for _, item := range items {
			feedmap := SiteRecord{
				title:      item.Title,
				URL:        item.Link,
				image:      getImageFromFeed(item.Content),
				updateDate: item.Published,
				siteID:     siteinfo.ID,
			}
			if feedmap.updateDate > siteinfo.latestDate {
				update_count++
				feedArray = append(feedArray, feedmap)
				if !isVisit[siteinfo.ID] {
					updateLatestDate(siteinfo.ID, feedmap.updateDate)
					isVisit[siteinfo.ID] = true
				}
			}
		}
	}
	_ = registerLatestArticleToDB(feedArray)
	// esId := registerLatestArticleToDB(feedArray)
	// registerLatestArticleToES(esId, feedArray)
	return update_count
}

func getImageFromFeed(feed string) string {
	reader := strings.NewReader(feed)
	doc, _ := goquery.NewDocumentFromReader(reader)
	imageUrl, _ := doc.Find("img").Attr("src")
	return imagectl.ArrangeImageUrl(imageUrl)
}

func registerLatestArticleToDB(articleList []SiteRecord) []int {
	db := openDB()
	defer db.Close()
	sql01_02 := "INSERT INTO /* sql01_02 */ articleTBL (title, URL, image, updateDate, click, siteID) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var esIdList []int
	var esId int
	for _, article := range articleList {
		err := db.QueryRow(sql01_02, article.title, article.URL, article.image, article.updateDate, 0, article.siteID).Scan(&esId)
		checkError(err)
		esIdList = append(esIdList, esId)
	}
	return esIdList
}

func jsonStruct(esIt int, doc SiteRecord) string {
	db := openDB()
	defer db.Close()
	sql01_03 := "SELECT title FROM siteTBL WHERE ID = $1"
	var sitetitle string
	err := db.QueryRow(sql01_03, doc.siteID).Scan(&sitetitle)
	docStruct := map[string]interface{}{
		"id":          esIt,
		"image":       doc.image,
		"publishedAt": doc.updateDate,
		"titles":      doc.title,
		"url":         doc.URL,
		"sitetitle":   sitetitle,
	}
	b, err := json.Marshal(docStruct)
	checkError(err)
	return string(b)
}

func registerLatestArticleToES(esIdList []int, articleList []SiteRecord) {
	jsonstrings := ""
	for i := 0; i < len(esIdList); i++ {
		jsonstrings += "{\"create\":{ \"_index\" : \"test_es\" , \"_id\" : \"" + strconv.Itoa(esIdList[i]) + "\"}}\n"
		jsonstrings += jsonStruct(esIdList[i], articleList[i]) + "\n"
	}

	print(jsonstrings)
	client := elastic.OpenES()

	res, err := client.Bulk(
		strings.NewReader(jsonstrings),
	)
	defer res.Body.Close()
	checkError(err)
}

func updateLatestDate(siteID int, updateDate string) {
	db := openDB()
	defer db.Close()

	sql01_03 := "UPDATE /* sql01_03 */ siteTBL SET latestDate = $1 WHERE ID = $2"
	stmt, err := db.Prepare(sql01_03)
	checkError(err)
	defer stmt.Close()
	_, err = stmt.Exec(updateDate, siteID)
	checkError(err)
}

func getSiteInfoList() []SiteInfo {
	db := openDB()
	defer db.Close()

	siteinfo := SiteInfo{}
	siteinfolist := []SiteInfo{}
	sql01_01 := "SELECT /* sql01_01 */ ID, title, rssURL, latestDate FROM siteTBL"

	selectSiteInfoList, err := db.Query(sql01_01)
	checkError(err)
	defer selectSiteInfoList.Close()
	for selectSiteInfoList.Next() {
		err := selectSiteInfoList.Scan(
			&siteinfo.ID,
			&siteinfo.title,
			&siteinfo.rssURL,
			&siteinfo.latestDate,
		)
		checkError(err)
		siteinfolist = append(siteinfolist, siteinfo)
	}
	return siteinfolist
}
