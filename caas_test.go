package main

import (
	_"encoding/json"
	_"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestHello(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	hello(response, request)

	if response.Code != http.StatusOK {
        t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
    }

    assert.Equal(t, "Hello world!", response.Body.String())

}

func TestCatReturnsJSON(t *testing.T) {
	request, _ := http.NewRequest("GET", "/cat", nil)
	response := httptest.NewRecorder()

	cat(response, request)

	ct := response.HeaderMap["Content-Type"][0]
	if !strings.EqualFold(ct, "application/json") {
		t.Fatalf("Content-Type does not equal 'application/json'")
	}
}

func TestCatpic(t *testing.T) {
	request, _ := http.NewRequest("GET", "/catpic", nil)
	response := httptest.NewRecorder()

	catpic(response, request)

	if response.Code != http.StatusOK {
        t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
    }

	assert.Contains(t, response.Body.String(), "img src=")

}
