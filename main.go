package main

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
)

func main() {
	router := gin.Default()

	router.GET("/ping", pong)
	router.POST("/resize", postFormValidation(), resizeHandler)

	router.Run(":8000")
}

func pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func resizeHandler(c *gin.Context) {
	url := c.PostForm("url")
	width, _ := strconv.ParseUint(c.PostForm("width"), 10, 64)
	height, _ := strconv.ParseUint(c.PostForm("height"), 10, 64)

	resizedImage, _ := imageUtil(url, uint(height), uint(width))
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resizedImage, nil)

	if err != nil {
		c.String(400, err.Error())
	}

	file, err := os.Create("tmp/resized_image.jpg")
	if err != nil {
		c.String(400, err.Error())
	}

	defer file.Close()

	jpeg.Encode(file, resizedImage, nil)
	b64string := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(200, gin.H{
		"results": b64string,
	})
}

func postFormValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err1 := strconv.ParseUint(c.PostForm("width"), 10, 64)
		_, err2 := strconv.ParseUint(c.PostForm("height"), 10, 64)

		switch {
		case err1 != nil:
			c.JSON(400, gin.H{
				"error": "width should be numeric",
			})
			c.Abort()
		case err2 != nil:
			c.JSON(400, gin.H{
				"error": "height should be numeric",
			})
			c.Abort()
		}
		c.Next()
	}
}

func imageUtil(url string, height uint, width uint) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imgResponse, err := jpeg.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	resizedImage := resize.Resize(width, height, imgResponse, resize.Lanczos3)
	return resizedImage, nil
}
