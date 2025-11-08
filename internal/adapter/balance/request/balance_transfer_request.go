package request

type BalanceTransferRequest struct {
	Amount      int    `json:"amount"`
	TargetEmail string `json:"targetEmail"`
}
