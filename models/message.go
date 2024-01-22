package models

type Message struct {
	Id          int 	`json:"Id"`
	UserId      string  `json:"UserId"`
	OrderItems  string  `json:"OrderItems"`
	PhoneNumber string  `json:"PhoneNumber"`
	TotalPrice  string  `json:"TotalPrice"`
}