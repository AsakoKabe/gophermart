package models

type Order struct {
	ID         string
	Num        string
	UserID     string
	UploadedAt string
}

type OrderWithAccrual struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}
