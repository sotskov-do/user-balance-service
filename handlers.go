package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var user User = User{Id: 1, Balance: 200} // TODO получить пользователя по id из БД

func operationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "wrong method")))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	fmt.Fprintln(w, user)

	params := &Operation{}
	err = json.Unmarshal(body, params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, params)

	// id := r.URL.Query().Get("id")
	// if id == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(fmt.Sprintf("{error: %v}", "query string dont have user id")))
	// 	return
	// }
	// fmt.Fprintln(w, id)

	if r.URL.Query().Get("a") == "2" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "Bad Request")))
		return
	}
	fmt.Fprintln(w, "operationHandler")
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Query())
	w.Write([]byte("!!!"))
	// fmt.Fprintln(w, r.Body)
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "wrong method")))
		return
	}

	fmt.Fprintln(w, "transferHandler")
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Query())
	w.Write([]byte("!!!"))
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "wrong method")))
		return
	}

	fmt.Fprintln(w, "balanceHandler")
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Query())
	w.Write([]byte("!!!"))
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("{error: %v}", "wrong method")))
		return
	}

	fmt.Fprintln(w, "historyHandler")
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Query())
	w.Write([]byte("!!!"))
}
