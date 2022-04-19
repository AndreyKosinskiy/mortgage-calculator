package app

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
)

// calcMMP function for calculation month mortgage payment
func calcMMP(resp *MortgageCalcResponse, b *models.Bank, name string, initLoan, downPayment float64) {
	validMessages := make([]string, 2)
	if initLoan > b.MaxLoan || downPayment < b.MinDownPayment {
		if initLoan > b.MaxLoan {
			validMessages[0] = "Initial loan not satisfies the maximum loan boundary of the bank."
		}
		if downPayment < b.MinDownPayment {
			validMessages[1] = "Down payment not satisfies the minimum down payment boundary of the bank."
		}
		if validMessages[0] != "" && validMessages[1] != "" {
			resp.ValidationMsg = strings.Join(validMessages, "\n")
		} else {
			resp.ValidationMsg = strings.Join(validMessages, "")
		}
	} else {
		//calculate monthly mortgage payment
		mmp := (float64(initLoan) * (b.Rate / 100 / 12) * math.Pow((1+b.Rate/100/12), float64(b.LoanTerm))) / (math.Pow((1+b.Rate/100/12), float64(b.LoanTerm)) - 1)
		resp.MonthMartgagePaymentMsg = fmt.Sprintf("Your month mortgage payment: %.2f$", mmp)
		resp.MortgageCalcRow = fmt.Sprintf("%s;%.2f$;%.2f$;%s;%.2f$", time.Now().Format("02.01.2006 15:04:05"), initLoan, downPayment, name, mmp)
	}
}
