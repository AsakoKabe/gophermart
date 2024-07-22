package models

type Withdrawal struct {
	ID          string  `json:"id,omitempty"`
	OrderNum    string  `json:"order"`
	Sum         float64 `json:"sum"`
	UserID      string  `json:"user_id,omitempty"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}
