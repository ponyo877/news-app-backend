package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	_ "os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo"
	"github.com/ponyo877/news-app-backend/handler/mongo"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	orgmongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

func UploadAvatar(file multipart.File, fileHeader *multipart.FileHeader, deviceHash string) string {
	bktName := "img.gitouhon-juku-k8s2.ga"
	imageBaseURL := "https://img.gitouhon-juku-k8s2.ga/"
	credentialFilePath := "config/config_gcp.json"

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	defer client.Close()

	body, err := ioutil.ReadAll(file)
	nowTimestamp := time.Now().Format("20060102150405")
	format := strings.Split(fileHeader.Filename, ".")[1]
	objName := "avatar/" + deviceHash + "_" + nowTimestamp + "." + format
	wc := client.Bucket(bktName).Object(objName).NewWriter(ctx)
	defer wc.Close()
	if _, err := wc.Write(body); err != nil {
		fmt.Println("createFile: unable to write data to bucket %q, file %q: %v", bktName, objName, err)
	}
	checkError(err)
	return imageBaseURL + objName
}

// FileUpload: https://echo.labstack.com/cookbook/file-upload/
func UpdateUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		qname := c.FormValue("name")
		qdeviceHash := c.FormValue("devicehash")
		avatar, avaterHeader, err := c.Request().FormFile("avatar")
		hasAvatar := true
		if err != nil {
			fmt.Println("Failed to open file at path")
			hasAvatar = false
		}
		fmt.Println(qname)

		ctx := context.Background()
		client := mongo.OpenMongo()
		err = client.Connect(ctx)
		defer client.Disconnect(ctx)
		checkError(err)

		usr_col := client.Database("newsdb").Collection("user_col")
		usr_filter := bson.M{
			"deviceHash": bson.M{"$eq": qdeviceHash},
		}

		usr_info := bson.M{}
		if qname != "" {
			usr_info["name"] = qname
		}
		if hasAvatar {
			qavatar := UploadAvatar(avatar, avaterHeader, qdeviceHash)
			usr_info["avatar"] = qavatar
		}

		usr_update := bson.M{
			"$set": usr_info,
		}
		usr_opts := options.Update().SetUpsert(true)
		_, err = usr_col.UpdateOne(ctx, usr_filter, usr_update, usr_opts)
		if err == orgmongo.ErrNoDocuments {
			return c.JSON(http.StatusOK, map[string]string{"Status": "NG"})
		}
		checkError(err)
		return c.JSON(http.StatusOK, map[string]string{"Status": "Ok"})
	}
}
