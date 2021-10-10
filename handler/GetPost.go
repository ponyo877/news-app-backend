package handler

import (
    "net/http"
    "strconv"
    "github.com/labstack/echo"
)

type feedRecord struct {
	ID         int
	title      string
	URL        string
	image      string
	updateDate string
	click      int
    siteID     int
    siteTitle  string
}

func GetPost() echo.HandlerFunc {
    return func(c echo.Context) error {		
		feed := feedRecord{}
        var feedArray []map[string]interface{}
        sql01_01 := "SELECT /* sql01_01 */ A.id, A.title, A.URL, A.image, A.updateDate, A.click, S.title FROM articleTBL A INNER JOIN siteTBL S ON A.siteID = S.ID ORDER BY updateDate DESC LIMIT 15"
        db := openDB()
        defer db.Close()
        selectFeedList, err := db.Query(sql01_01)
        checkError(err)      	
        defer selectFeedList.Close()
        for selectFeedList.Next() {
        	err = selectFeedList.Scan(
        		&feed.ID,
        		&feed.title,
        		&feed.URL,
                &feed.image,
                &feed.updateDate,
                &feed.click,
                // &feed.siteID,
                &feed.siteTitle,
            )
            checkError(err) 
            feedmap := map[string]interface{}{
                "id":          feed.ID,
                "titles":      feed.title,
                "url":         feed.URL,
                "image":       feed.image,
                "publishedAt": feed.updateDate,
                "sitetitle":   feed.siteTitle,
            }
            feedArray = append(feedArray, feedmap)
        }
        return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
    }
}

func GetPostFromTo() echo.HandlerFunc {
    return func(c echo.Context) error {		
        qfrom, _ := strconv.Atoi(c.QueryParam("from"))
        // qto, _ := strconv.Atoi(c.QueryParam("to"))
        // qlimit := qto - qfrom
		feed := feedRecord{}
        var feedArray []map[string]interface{}
        sql01_02 := "SELECT /* sql01_02 */ A.id, A.title, A.URL, A.image, A.updateDate, A.click, S.title FROM articleTBL A INNER JOIN siteTBL S ON A.siteID = S.ID ORDER BY updateDate DESC LIMIT 15 OFFSET $1"
        db := openDB()
        defer db.Close()
        selectFeedList, err := db.Query(sql01_02, strconv.Itoa(qfrom))
        checkError(err)      	
        defer selectFeedList.Close()
        for selectFeedList.Next() {
        	err = selectFeedList.Scan(
        		&feed.ID,
        		&feed.title,
        		&feed.URL,
                &feed.image,
                &feed.updateDate,
                &feed.click,
                // &feed.siteID,
                &feed.siteTitle,
            )
            checkError(err) 
            feedmap := map[string]interface{}{
                "id":          feed.ID,
                "titles":      feed.title,
                "url":         feed.URL,
                "image":       feed.image,
                "publishedAt": feed.updateDate,
                "sitetitle":   feed.siteTitle,
            }
            feedArray = append(feedArray, feedmap)
        }
        return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": feedArray})
    }
}