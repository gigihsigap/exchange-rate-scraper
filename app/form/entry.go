package form

type buysell struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

type Entry struct {
	Date   string  `json:"date"`
	Symbol string  `json:"symbol"`
	ER     buysell `json:"e_rate"`
	TT     buysell `json:"tt_counter"`
	BN     buysell `json:"bank_notes"`
}
