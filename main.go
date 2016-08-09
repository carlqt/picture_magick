package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nfnt/resize"
	"net/http"
	_ "strconv"
)

func main() {
	router := gin.Default()

	router.GET("/ping", pong)
	router.POST

	router.Run(":8000")
}

func pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
