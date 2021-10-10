package handler

import (
	"net/http"
	"./elastic"
    "github.com/labstack/echo"
)

func GetSearch() echo.HandlerFunc {
    return func(c echo.Context) error {	
		searchwords := c.QueryParam("words")
		searchResult := elastic.GetSearchResultTmp(searchwords)
        return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": searchResult})
    }
}
