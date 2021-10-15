package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ponyo877/news-app-backend/handler/elastic"
)

func GetSearch() echo.HandlerFunc {
	return func(c echo.Context) error {
		searchwords := c.QueryParam("words")
		searchResult := elastic.GetSearchResultTmp(searchwords)
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{"data": searchResult})
	}
}
