package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
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
