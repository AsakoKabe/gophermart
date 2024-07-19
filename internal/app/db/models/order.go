package models

type Order struct {
	ID         string
	Num        int
	UserID     string
	UploadedAt string
}

type OrderWithAccrual struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
