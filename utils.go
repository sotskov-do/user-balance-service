package main

import (
	"fmt"
	"net/http"
)

func idChecker(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "enter user id")))
		return
	}
}
