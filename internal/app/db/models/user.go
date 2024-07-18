package models

type User struct {
	ID       string
	Login    string `json:"login"`
	Password string `json:"password"`
}
