package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func operationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applciation/json")
	if r.Method != "POST" {
		http.Error(w, "{\"error\": \"wrong request method\"}", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err) // TODO обработать
	}
	defer r.Body.Close()

	params := &Operation{}
	err = json.Unmarshal(body, params)
	if err != nil {
		http.Error(w, "{\"error\": \"provide correct id, type and amount in request body\"}", http.StatusBadRequest)
		return
	}

	if validErrs := params.validate(); len(validErrs) > 0 {
		errs := map[string]interface{}{"error": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	user, err := getUser(params.Id)
	if err != nil && params.Type == "credit" {
		user, err = createUser(params.Id, 0)
		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
			return
		}
	}

	log.Println(user, params)

	if params.Type == "credit" {
		user.Balance += params.Amount
	} else {
		user.Balance -= params.Amount
	}

	if user.Balance < 0 {
		http.Error(w, "{\"error\": \"not enough money\"}", http.StatusBadRequest)
		return
	}

	err = updateUserBalance(user.Id, user.Balance)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}

	// TODO записать операцию в историю

	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"balance\": %v}", user.Balance)))
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applciation/json")
	if r.Method != "POST" {
		http.Error(w, "{\"error\": \"wrong request method\"}", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err) // TODO обработать
	}
	defer r.Body.Close()

	params := &Transfer{}
	err = json.Unmarshal(body, params)
	if err != nil {
		http.Error(w, "{\"error\": \"provide correct sender_id, reciever_id and amount in request body\"}", http.StatusBadRequest)
		return
	}

	if validErrs := params.validate(); len(validErrs) > 0 {
		errs := map[string]interface{}{"error": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	user_sender, err := getUser(params.SenderId)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}
	user_reciever, err := getUser(params.RecieverId)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}

	log.Println(user_sender, user_reciever, params)

	user_sender.Balance -= params.Amount
	user_reciever.Balance += params.Amount

	if user_sender.Balance < 0 {
		http.Error(w, "{\"error\": \"not enough money\"}", http.StatusBadRequest)
		return
	}

	err = updateUserBalance(user_sender.Id, user_sender.Balance)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}
	err = updateUserBalance(user_reciever.Id, user_reciever.Balance)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}

	// TODO записать операцию в историю

	fmt.Fprintln(w, "transferHandler")
	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"sender_balance\": %v, \"reciever_balance\": %v}", user_sender.Balance, user_reciever.Balance)))
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	var balance float64

	w.Header().Set("Content-type", "applciation/json")
	if r.Method != "GET" {
		http.Error(w, "{\"error\": \"wrong request method\"}", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("id") == "" {
		http.Error(w, "{\"error\": \"add user id to query string\"}", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id")) // TODO обработать ошибку
	log.Println(id)                                // TODO delete

	user, err := getUser(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
		return
	}

	currency := strings.ToUpper(r.URL.Query().Get("currency"))
	if currency == "" || currency == "RUB" {
		currency = "RUB"
		balance = user.Balance
	} else {
		var err error
		balance, err = currencyConversion(currency, user.Balance)
		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
			return
		}
	}

	w.Write([]byte(fmt.Sprintf("{\"id\": %v, \"currency\": %v, \"balance\": %v}", user.Id, currency, balance)))
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applciation/json")
	if r.Method != "GET" {
		http.Error(w, "{\"error\": \"wrong request method\"}", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("id") == "" {
		http.Error(w, "{\"error\": \"add user id to query string\"}", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	fmt.Println(id)

	// TODO получить пользователя по id из БД
	// TODO проверка на существование пользователя
	// TODO вывести операции (таблица вида - user_id, тип операции (debit, credit, transfer), сумма, дата операции)

	fmt.Fprintln(w, "historyHandler")
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Query())
	w.Write([]byte("!!!"))
}
