package transaction

type Transaction struct {
	ID     int     `json:"Id"`
	Date   string  `json:"Date"`
	Amount float64 `json:"Amount"`
}
