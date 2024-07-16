package models

import "errors"

type User struct {
	ID       string
	Login    string `json:"login"`
	Password string `json:"password"`
}

var LoginAlreadyExist = errors.New("user login already exist")
