package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var user1 User = User{Id: 1, Balance: 200}
var user2 User = User{Id: 2, Balance: 100}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-type", "applciation/json")
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

	// TODO получить пользователя по id из БД
	// TODO создать пользователя при зачислении, если не найден
	fmt.Println(user1, params)

	if params.Type == "credit" {
		user1.Balance += params.Amount
	} else {
		user1.Balance -= params.Amount
	}

	if user1.Balance < 0 {
		http.Error(w, "{\"error\": \"not enough money\"}", http.StatusBadRequest)
		return
	}

	// TODO записать результат в базу
	// TODO записать операцию в историю

	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"balance\": %v}", user1.Balance)))
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
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

	// TODO получить пользователей по id из БД
	fmt.Println(user1, user2, params)

	user2.Balance += params.Amount
	user1.Balance -= params.Amount

	if user1.Balance < 0 {
		http.Error(w, "{\"error\": \"not enough money\"}", http.StatusBadRequest)
		return
	}

	// TODO записать результат в базу
	// TODO записать операцию в историю

	fmt.Fprintln(w, "transferHandler")
	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"sender_balance\": %v, \"reciever_balance\": %v}", user1.Balance, user2.Balance)))
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	var balance float64

	if r.Method != "GET" {
		http.Error(w, "{\"error\": \"wrong request method\"}", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("id") == "" {
		http.Error(w, "{\"error\": \"add user id to query string\"}", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	fmt.Println(id) // TODO delete

	currency := strings.ToUpper(r.URL.Query().Get("currency"))
	if currency == "" || currency == "RUB" {
		currency = "RUB"
		balance = user1.Balance
	} else {
		var err error
		balance, err = currencyConversion(currency, user1.Balance)
		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusBadRequest)
			return
		}
	}

	// TODO получить пользователя по id из БД
	// TODO проверка на существование пользователя
	w.Write([]byte(fmt.Sprintf("{\"id\": %v, \"currency\": %v, \"balance\": %v}", user1.Id, currency, balance)))
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
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
