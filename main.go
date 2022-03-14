package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/operation/", operationHandler)
	http.HandleFunc("/transfer/", transferHandler)
	http.HandleFunc("/balance/", balanceHandler)
	http.HandleFunc("/history/", historyHandler)

	http.ListenAndServe(":8080", nil)
}
