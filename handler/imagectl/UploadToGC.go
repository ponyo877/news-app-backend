package imagectl

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"bytes"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
	"hash/fnv"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
	"strconv"
	"strings"
)

func HashID(imageUrl string) uint32 {
	h := fnv.New32()
	h.Write([]byte(imageUrl))
	sum := h.Sum32()
	return sum
}

func getRandomImage() string {
	// http://flat-icon-design.com/
	RandomImageBaseURL := "https://img.gitouhon-juku-k8s2.ga/default_"
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNum := strconv.Itoa(r1.Intn(10))
	return RandomImageBaseURL + randomNum + ".jpg"
}

func ArrangeImageUrl(imageUrl string) string {
	bktName := "img.gitouhon-juku-k8s2.ga"
	imageBaseURL := "https://img.gitouhon-juku-k8s2.ga/"
	credentialFilePath := "config_gcp.json"
	if imageUrl == 	"" || !strings.HasPrefix(imageUrl, "http") {
		return getRandomImage()
	}
	// fmt.Println(imageUrl) // for imageUrl debug 
	response, err := http.Get(imageUrl)
	defer response.Body.Close()
	checkError(err)

	if response.StatusCode == 200 {
		ctx := context.Background()
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
		defer client.Close()

		body, err := ioutil.ReadAll(response.Body)
		r := bytes.NewReader(body)
		_, format, err := image.DecodeConfig(r)
		checkError(err)
		if format != "jpeg" && format != "png" {
			return getRandomImage()
		}
		objName := fmt.Sprint(HashID(imageUrl)) + "." + format
		wc := client.Bucket(bktName).Object(objName).NewWriter(ctx)
		defer wc.Close()
		if _, err := wc.Write(body); err != nil {
			// fmt.Println("createFile: unable to write data to bucket %q, file %q: %v", bktName, objName, err)
			return getRandomImage()
		}
		return imageBaseURL + objName
	} else {
		return getRandomImage()
	}
}

// checkError(err)
func UploadToGC() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bktName := "gitouhon-juku-k8s2.ga"
		bktName := "img.gitouhon-juku-k8s2.ga"
		objName := "sample_cat.jpg"
		imageUrl := "http://placekitten.com/g/640/340"
		credentialFilePath := "config_gcp.json"
		ctx := context.Background()

		client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
		defer client.Close()

		response, err := http.Get(imageUrl)
		checkError(err)
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		checkError(err)

		wc := client.Bucket(bktName).Object(objName).NewWriter(ctx)
		defer wc.Close()
		if _, err := wc.Write(body); err != nil {
			fmt.Println("createFile: unable to write data to bucket %q, file %q: %v", bktName, objName, err)
			return c.JSON(http.StatusOK, map[string]string{"status": "NG"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
