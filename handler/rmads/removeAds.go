package rmads

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func RemoveAds(Url string) (modifiedHtml string) {
	httpres := getHttpResponse(Url)
	// arrange header
	
}

func getHttpResponse(Url string) string {
	req, _ := http.NewRequest("GET", Url, nil)
	userAgent := "Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36"
	req.Header.Set("User-Agent", userAgent)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http client error occur")
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray)
}