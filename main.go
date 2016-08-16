package main

import (
	"os"
	"fmt"
	"github.com/urfave/cli"
	"bytes"
	// "github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	// "os"
	b64 "encoding/base64"
	"io"
	"strconv"
)

func main() {
	var port string

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "port, p",
			Value: "8000",
			Usage: "Run server in a specific port",
		},
	}

	app.Name = "Picture Magick"
	app.Usage = "Resizing images"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		port = fmt.Sprintf(":%s", c.String("port"))
		// run the router
		router := gin.Default()

		router.Use(corsHeader)
		router.GET("/ping", Pong)
		router.POST("/resize", PostFormValidation(), ResizeHandler)

		router.Run(port)
		return nil
	}

	app.Run(os.Args)

}

func Pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// This method is here to help in writing the test
func Foo() string {
	return "hello"
}

func ResizeHandler(c *gin.Context) {
	var buf io.Writer

	url := c.PostForm("url")
	width, _ := strconv.ParseUint(c.PostForm("width"), 10, 64)
	height, _ := strconv.ParseUint(c.PostForm("height"), 10, 64)

	img, imgType, err := imageUtil(url)

	if err != nil {
		c.String(400, err.Error())
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	switch imgType {
	case "jpeg":
		buf = jpegEncode(resizedImage)
	case "png":
		buf = pngEncode(resizedImage)
	case "gif":
		buf = gifEncode(resizedImage)
	}

	c.JSON(200, gin.H{
		"results": "Image successfully converted",
		"image":   encodeBase64(buf),
	})
}

func PostFormValidation() gin.HandlerFunc {
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
}

func jpegEncode(img image.Image) io.Writer {
	// file, _ := os.Create("tmp/resized_image.jpg")
	newBuff := bytes.NewBuffer(nil)

	// defer file.Close()
	jpeg.Encode(newBuff, img, nil)
	return newBuff
}

func pngEncode(img image.Image) io.Writer {
	// file, err := os.Create("tmp/resized_image.png")
	newBuff := bytes.NewBuffer(nil)

	// defer file.Close()
	png.Encode(newBuff, img)
	return newBuff
}

func gifEncode(img image.Image) io.Writer {
	// file, err := os.Create("tmp/resized_image.gif")
	newBuff := bytes.NewBuffer(nil)

	// defer file.Close()
	gif.Encode(newBuff, img, nil)
	return newBuff
}

func encodeBase64(buf io.Writer) string {
	bufBytes := buf.(*bytes.Buffer)
	return b64.StdEncoding.EncodeToString(bufBytes.Bytes())
}

// CORS Header
func corsHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
