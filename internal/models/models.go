package models

// Bank ...
type Bank struct {
	Name           string
	Rate           float64
	MaxLoan        uint
	MinDownPayment uint
	LoanTerm       uint
}

// Bank ...
type MortgageCalc struct {
	InitLoan    uint
	DownPayment uint
	BankName    string
}
