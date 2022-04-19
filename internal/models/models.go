package models

// Bank ...
type Bank struct {
	Name           string
	Rate           float64
	MaxLoan        float64
	MinDownPayment float64
	LoanTerm       uint
}

// Bank ...
type MortgageCalc struct {
	InitLoan    uint
	DownPayment uint
	BankName    string
}
