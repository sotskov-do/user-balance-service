package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Operation Tests
func TestOperationComplete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"id\": 1, \"type\": \"credit\", \"amount\": 100}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["result"] != "success" {
		t.Fail()
	}
}

func TestOperationCreateNewUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"id\": 2, \"type\": \"credit\", \"amount\": 100}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["result"] != "success" {
		t.Fail()
	}
}

func TestOperationNotEnoughMoney(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"id\": 1, \"type\": \"debit\", \"amount\": 1000000000}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "not enough money" {
		t.Fail()
	}
}

func TestOperationEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "provide correct id, type and amount in request body" {
		t.Fail()
	}
}

func TestOperationBodyError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"id\": -1, \"type\": \"crefdit\", \"amount\": -1}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["id"].([]interface{})[0] != "id must be positive integer more than zero" {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["type"].([]interface{})[0] != "type must be debit or credit" {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["amount"].([]interface{})[0] != "amount must be positive integer more than zero" {
		t.Fail()
	}
}

func TestOperationWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Get(ts.URL + "/operation/")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "wrong request method" {
		t.Fail()
	}
}

func TestOperationNoSuchId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(operationHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"id\": 1000000000, \"type\": \"debit\", \"amount\": 1}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "no user with such id" {
		t.Fail()
	}
}

// Transfer Tests
func TestTransferComplete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/transfer/", "application/json", strings.NewReader("{\"sender_id\": 1, \"reciever_id\": 2, \"amount\": 1}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["result"] != "success" {
		t.Fail()
	}
}

func TestTransferWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Get(ts.URL + "/transfer/")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "wrong request method" {
		t.Fail()
	}
}

func TestTransferEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/transfer/", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "provide correct sender_id, reciever_id and amount in request body" {
		t.Fail()
	}
}

func TestTransferErrorBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/transfer/", "application/json", strings.NewReader("{\"asender_id\": 1, \"areciever_id\": 2, \"aamount\": 1}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["sender_id"].([]interface{})[0] != "sender_id must be positive integer more than zero" {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["reciever_id"].([]interface{})[0] != "reciever_id must be positive integer more than zero" {
		t.Fail()
	}
	if result["errors"].(map[string]interface{})["amount"].([]interface{})[0] != "amount must be positive integer more than zero" {
		t.Fail()
	}
}

func TestTransferNotEnoughMoney(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"sender_id\": 1, \"reciever_id\": 2, \"amount\": 1000000000}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "not enough money" {
		t.Fail()
	}
}
func TestTransferSameId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"sender_id\": 1, \"reciever_id\": 1, \"amount\": 1}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "sender and receiver are the same" {
		t.Fail()
	}
}

func TestTransferNoSuchSenderId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"sender_id\": 1000000000, \"reciever_id\": 2, \"amount\": 1000000000}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "no user with such id" {
		t.Fail()
	}
}

func TestTransferNoSuchRecieverId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(transferHandler))
	defer ts.Close()
	res, err := http.Post(ts.URL+"/operation/", "application/json", strings.NewReader("{\"sender_id\": 1, \"reciever_id\": 1000000000, \"amount\": 1000000000}"))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "no user with such id" {
		t.Fail()
	}
}

// Balance Tests
func TestBalanceComplete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/balance/?id=1")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	_, ok := result["balance"]
	if !ok {
		t.Fail()
	}
}

func TestBalanceWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Post(ts.URL+"/balance/?id=1", "", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "wrong request method" {
		t.Fail()
	}
}

func TestBalanceNoId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/balance/")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "add user id to query string" {
		t.Fail()
	}
}

func TestBalanceNoSuchId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/balance/?id=1000000000")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "no user with such id" {
		t.Fail()
	}
}

func TestBalanceDifferentCurrency(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/balance/?id=1&currency=usd")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["currency"] != "USD" {
		t.Fail()
	}
}

func TestBalanceWrongCurrency(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(balanceHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/balance/?id=1&currency=abc")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "can't find such currency. please use another one" {
		t.Fail()
	}
}

// History Tests
func TestHistoryComplete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/?id=1")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := []map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	_, ok := result[0]["datetime"]
	if !ok {
		t.Fail()
	}
}

func TestHistoryWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Post(ts.URL+"/history/?id=1", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "wrong request method" {
		t.Fail()
	}
}

func TestHistoryNoId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "add user id to query string" {
		t.Fail()
	}
}

func TestHistoryNoSuchId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/?id=1000000")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "no user with such id" {
		t.Fail()
	}
}

func TestHistoryWrongSorted(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/?id=1&sorted=abc")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "sorted must be amount or datetime" {
		t.Fail()
	}
}

func TestHistoryWrongOrder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/?id=1&order=abc")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "order must be asc or desc" {
		t.Fail()
	}
}

func TestHistoryNoHistory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(historyHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/history/?id=1&page=1000000")
	if err != nil {
		t.Fail()
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fail()
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fail()
	}
	if result["error"] != "nothing found" {
		t.Fail()
	}
}
