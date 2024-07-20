package models

type User struct {
	ID       string `json:"ID,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
