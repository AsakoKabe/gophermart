package models

type User struct {
	ID         string  `json:"ID,omitempty"`
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	Accruals   float64 `json:"accruals,omitempty"`
	Withdrawal float64 `json:"withdrawal,omitempty"`
}
