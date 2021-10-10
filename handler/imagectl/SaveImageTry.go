package imagectl

import (
	"io/ioutil"
	"net/http"
	"os"
	_"context"
	"github.com/labstack/echo"
)

func SaveImageToCS() echo.HandlerFunc{
	return func(c echo.Context) error {	
		imageUrl := "http://placekitten.com/g/640/340"
		response, err := http.Get(imageUrl)
		checkError(err)
		defer response.Body.Close()
		
		body, err := ioutil.ReadAll(response.Body)
		checkError(err)

		file, err := os.Create("sample_cat.jpg")
		checkError(err)

		defer file.Close()
		file.Write(body)
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
