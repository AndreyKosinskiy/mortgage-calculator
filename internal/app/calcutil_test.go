package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
)

func TestCalcMMP(t *testing.T) {
	type calcInputs struct {
		initLoan    float64
		downPayment float64
		name        string
	}
	b := &models.Bank{Name: "TestBank", Rate: 10, MaxLoan: 1000, MinDownPayment: 10, LoanTerm: 12}
	table := []struct {
		caseName string
		b        *models.Bank
		inputs   calcInputs
		want     MortgageCalcResponse
		got      string
	}{
		{
			caseName: "succes calculation",
			b:        b,
			inputs:   calcInputs{500, 15, "TestBank"},
			want: MortgageCalcResponse{
				BankNameList:            []string{b.Name},
				ValidationMsg:           "",
				MonthMartgagePaymentMsg: "Your month mortgage payment: 43.96$",
				MortgageCalcRow:         fmt.Sprintf("%s;%.2f$;%.2f$;%s;%.2f$", time.Now().Format("02.01.2006 15:04:05"), 500.0, 15.0, b.Name, 43.96),
			},
		},
		{
			caseName: "down payment validation",
			b:        b,
			inputs:   calcInputs{500, 5, "TestBank"},
			want: MortgageCalcResponse{
				BankNameList:            []string{b.Name},
				ValidationMsg:           "Down payment not satisfies the minimum down payment boundary of the bank.",
				MonthMartgagePaymentMsg: "",
				MortgageCalcRow:         "",
			},
		},
		{
			caseName: "max loan validation",
			b:        b,
			inputs:   calcInputs{5000, 15, "TestBank"},
			want: MortgageCalcResponse{
				BankNameList:            []string{b.Name},
				ValidationMsg:           "Initial loan not satisfies the maximum loan boundary of the bank.",
				MonthMartgagePaymentMsg: "",
				MortgageCalcRow:         "",
			},
		},
		{
			caseName: "max loan and down paymnet validation",
			b:        b,
			inputs:   calcInputs{5000, 5, "TestBank"},
			want: MortgageCalcResponse{
				BankNameList:            []string{b.Name},
				ValidationMsg:           "Initial loan not satisfies the maximum loan boundary of the bank.\nDown payment not satisfies the minimum down payment boundary of the bank.",
				MonthMartgagePaymentMsg: "",
				MortgageCalcRow:         "",
			},
		},
	}
	for _, tc := range table {
		t.Run(tc.caseName, func(t *testing.T) {
			got := MortgageCalcResponse{BankNameList: []string{tc.b.Name}}
			fmt.Println(tc.b.Name)
			calcMMP(&got, tc.b, tc.b.Name, tc.inputs.initLoan, tc.inputs.downPayment)
			want := tc.want
			if got.MonthMartgagePaymentMsg != want.MonthMartgagePaymentMsg || got.MortgageCalcRow != want.MortgageCalcRow || got.ValidationMsg != want.ValidationMsg {
				t.Error(tc.caseName+": \nwant: ", want, " \nBut got: ", got)
			}
		})
	}
}
