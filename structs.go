package main

// type ServiceError struct {
// 	Error string
// }

type Operation struct {
	Id     int
	Type   string
	Amount int
}

type Transfer struct {
	SenderId   int
	RecieverId int
	Amount     int
}

type User struct {
	Id      int
	Balance int
}
