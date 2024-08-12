package handlers

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type withdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
