package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBalance(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/balance/?id=1000", nil)
	w := httptest.NewRecorder()
	balanceHandler(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	log.Println(string(data), err)
}

// TODO тесты
