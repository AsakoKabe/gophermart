package models

type Order struct {
	ID         string      `json:"-"`
	Num        string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual"`
	UserID     string      `json:"-"`
	UploadedAt string      `json:"uploaded_at,omitempty"`
}

type OrderStatus string

const NEW OrderStatus = "NEW"
const PROCESSED OrderStatus = "PROCESSED"
const INVALID OrderStatus = "INVALID"
const PROCESSING OrderStatus = "PROCESSING"
const REGISTERED OrderStatus = "REGISTERED"
