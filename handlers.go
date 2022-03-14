package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func operationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("{\"error\": \"wrong request method\"}"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	defer r.Body.Close()

	params := &Operation{}
	err = json.Unmarshal(body, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"provide correct id, type and amount in request body\"}"))
		return
	}

	if validErrs := params.validate(); len(validErrs) > 0 {
		errs := map[string]interface{}{"errors": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	user, err := getUser(params.Id)
	if err != nil && params.Type == "credit" {
		user, err = createUser(params.Id, 0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	if params.Type == "credit" {
		user.Balance += params.Amount
	} else {
		user.Balance -= params.Amount
	}

	if user.Balance < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"not enough money\"}"))
		return
	}

	err = updateUserBalance(user.Id, user.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	err = updateHistory(user.Id, params.Type, params.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"balance\": %v}", user.Balance)))
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("{\"error\": \"wrong request method\"}"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	defer r.Body.Close()

	params := &Transfer{}
	err = json.Unmarshal(body, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"provide correct sender_id, reciever_id and amount in request body\"}"))
		return
	}

	if validErrs := params.validate(); len(validErrs) > 0 {
		errs := map[string]interface{}{"errors": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	if params.SenderId == params.RecieverId {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"sender and receiver are the same\"}"))
		return
	}

	user_sender, err := getUser(params.SenderId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	user_reciever, err := getUser(params.RecieverId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	user_sender.Balance -= params.Amount
	user_reciever.Balance += params.Amount

	if user_sender.Balance < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"not enough money\"}"))
		return
	}

	err = updateUserBalance(user_sender.Id, user_sender.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	err = updateUserBalance(user_reciever.Id, user_reciever.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	err = updateHistory(user_sender.Id, "transfer_send", params.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	err = updateHistory(user_reciever.Id, "transfer_recieve", params.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"result\": \"success\", \"sender_balance\": %v, \"reciever_balance\": %v}", user_sender.Balance, user_reciever.Balance)))
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var balance float64

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("{\"error\": \"wrong request method\"}"))
		return
	}

	if r.URL.Query().Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"add user id to query string\"}"))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	user, err := getUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}
	}

	w.Write([]byte(fmt.Sprintf("{\"id\": %v, \"currency\": \"%v\", \"balance\": %v}", user.Id, currency, balance)))
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("{\"error\": \"wrong request method\"}"))
		return
	}

	if r.URL.Query().Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"add user id to query string\"}"))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}
	_, err = getUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	var page int
	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
			return
		}
	}

	var sorted string
	if r.URL.Query().Get("sorted") == "" {
		sorted = "datetime"
	} else {
		sorted = r.URL.Query().Get("sorted")
		if sorted != "amount" && sorted != "datetime" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\": \"sorted must be amount or datetime\"}"))
			return
		}
	}

	var order string
	if r.URL.Query().Get("order") == "" {
		order = "ASC"
	} else {
		order = strings.ToUpper(r.URL.Query().Get("order"))
		if order != "ASC" && order != "DESC" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\": \"order must be asc or desc\"}"))
			return
		}
	}

	historyPage, err := getHistory(id, page, sorted, order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	if len(historyPage) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"error\": \"nothing found\"}"))
		return
	} else {
		json.NewEncoder(w).Encode(historyPage)
	}
}
