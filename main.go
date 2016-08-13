package main

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
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

	img, imgType, err := imageUtil(url)
	// img, err := jpegDecode(url)
	// imgType := "jpeg"

	if err != nil {
		c.String(400, err.Error())
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	switch imgType {
	case "jpeg":
		jpegEncode(resizedImage)
	case "png":
		pngEncode(resizedImage)
	}


	c.JSON(200, gin.H{
		"results": "Image successfully converted",
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

func imageUtil(url string) (image.Image, string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	return image.Decode(response.Body)

	// resizedImage := resize.Resize(width, height, imgResponse, resize.Lanczos3)
	// return resizedImage, nil
}

func jpegDecode(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return jpeg.Decode(response.Body)
}

func jpegEncode(img image.Image) {
	file, _ := os.Create("tmp/resized_image.jpg")

	defer file.Close()
	jpeg.Encode(file, img, nil)
}

func pngEncode(img image.Image) {
	file, err := os.Create("tmp/resized_image.png")
	if err != nil {
		color.Red(err.Error())
	}

	defer file.Close()
	png.Encode(file, img)
}