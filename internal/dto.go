package internal

type Event struct {
	UserID      string `json:"userId"`
	From        string `json:"from"`
	To          string `json:"to"`
	AmountWei   string `json:"amountWei"`
	TxHash      string `json:"hash"`
	BlockNumber uint64 `json:"blockNumber"`
}
