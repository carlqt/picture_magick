package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
  "strings"
)

func TestPong(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/ping", Pong)
	r.ServeHTTP(w, req)

	res := w.Body.String()
	if res != "pong" {
		t.Errorf("Response should be pong but was %s", res)
	}
}

func TestFoo(t *testing.T) {
	foo := Foo()
	if foo != "hello" {
		t.Errorf("Return string should be `hello` but was %s", foo)
	}
}

func TestResizeHandler(t *testing.T) {
	data := url.Values{}
	data.Add("url", "http://www.vcahospitals.com/img/common/dog-care/breeds/pembroke-welsh-corgi.png")
	data.Add("width", "200")
	data.Add("height", "200")

	req, _ := http.NewRequest("POST", "/resize", strings.NewReader(data.Encode()))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/resize", ResizeHandler)
	r.ServeHTTP(w, req)

	res := w.Body.String()
	if res != "pong" {
		t.Errorf("Response should be pong but was %s", res)
	}
}
