package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/operation/", operationHandler)
	http.HandleFunc("/transfer/", transferHandler)
	http.HandleFunc("/balance/", balanceHandler)
	http.HandleFunc("/history/", historyHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
