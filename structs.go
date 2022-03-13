package main

import "net/url"

type Operation struct {
	Id     int
	Type   string
	Amount float64
}

func (o *Operation) validate() url.Values {
	var availableTypes map[string]bool = map[string]bool{"debit": true, "credit": true}
	errs := url.Values{}

	if o.Id <= 0 {
		errs.Add("id", "id must be positive integer more than zero")
	}

	_, ok := availableTypes[o.Type]
	if !ok {
		errs.Add("type", "type must be debit or credit")
	}

	if o.Amount <= 0 {
		errs.Add("amount", "amount must be positive integer more than zero")
	}

	return errs
}

type Transfer struct {
	SenderId   int `json:"sender_id"`
	RecieverId int `json:"reciever_id"`
	Amount     float64
}

func (t *Transfer) validate() url.Values {
	errs := url.Values{}

	if t.SenderId <= 0 {
		errs.Add("sender_id", "sender_id must be positive integer more than zero")
	}

	if t.RecieverId <= 0 {
		errs.Add("reciever_id", "reciever_id must be positive integer more than zero")
	}

	if t.Amount <= 0 {
		errs.Add("amount", "amount must be positive integer more than zero")
	}

	return errs
}

type User struct {
	Id      int
	Balance float64
}

type History struct {
	UserId   int     `json:"user_id"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Datetime string  `json:"datetime"`
	Idx      int     `json:"idx"`
}
