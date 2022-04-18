package app

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
	bankrepository "github.com/AndreyKosinskiy/mortgage-calculator/internal/repository/bankRepository"
)

// MortgageCalcResponse
type MortgageCalcResponse struct {
	BankNameList            []string
	ValidationMsg           string
	MonthMartgagePaymentMsg string
}

// BankListResponse
type BankListResponse struct {
	BankList []*models.Bank
}

// BankResponse
type BankResponse struct {
	Bank *models.Bank
}

// MortgageCalcHandler
func (a *App) MortgageCalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	tmpl, err := template.ParseFiles("../web/templates/index.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		a.logger.Printf("Error MortgageCalcHandler.ParseFiles: %s", err)
		return
	}

	resp := MortgageCalcResponse{}
	repo := bankrepository.New(a.db, *a.logger)
	a.logger.Printf("%v: %s %s %s ", time.Now(), "MortgageCalcHandler", r.URL.Path, r.Method)

	bs, err := repo.BankList(r.Context())
	if err != nil {
		http.Error(w, "can`t get bank list", 404)
		a.logger.Printf("Error MortgageCalcHandler.BankList: %s", err)
		return
	}
	bn := make([]string, len(bs))
	for i, v := range bs {
		bn[i] = v.Name
	}
	resp.BankNameList = bn

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "can`t parse Form from bank_create.html", 404)
			a.logger.Printf("Error MortgageCalcHandler.ParseForm: %s", err)
			return
		}

		initLoan, err := strconv.ParseFloat(r.PostFormValue("init-loan"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[init-loan]", 404)
			a.logger.Printf("Error MortgageCalcHandler.PostFormValue.init-loan: %s", err)
			return
		}
		downPayment, err := strconv.ParseFloat(r.PostFormValue("down-payment"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[down-payment]", 404)
			a.logger.Printf("Error MortgageCalcHandler.PostFormValue.down-payment: %s", err)
			return
		}
		name := r.PostFormValue("name")

		b, err := repo.BankByName(r.Context(), name)
		if err != nil {
			http.Error(w, "can`t get bank by name", 404)
			a.logger.Printf("Error MortgageCalcHandler.BankByName: %s", err)
			return
		}

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
		}
	}
	err = tmpl.Execute(w, resp)
	if err != nil {

		http.Error(w, "can`t execute template", 404)
	}
}

func (a *App) BankListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bank-list" {
		http.Error(w, "can`t parse form index.html", 404)
		return
	}
	tmpl, err := template.ParseFiles("../web/templates/bank_list.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	resp := BankListResponse{}
	repo := bankrepository.New(a.db, *a.logger)
	a.logger.Printf("%v: %s %s %s ", time.Now(), "BankListHandler", r.URL.Path, r.Method)

	switch r.Method {
	case http.MethodGet:
		bs, err := repo.BankList(r.Context())
		if err != nil {
			http.Error(w, "can`t parse template index.html", 404)
			return
		}
		resp.BankList = bs
		a.logger.Println("bs len: ", len(bs))
		tmpl.Execute(w, resp)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "can`t parse Form from bank_create.html", 404)
			return
		}
		name := r.PostFormValue("name")
		rate, err := strconv.ParseFloat(r.PostFormValue("rate"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[rate]", 404)
		}
		maxLoan, err := strconv.ParseFloat(r.PostFormValue("max-loan"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[max-loan]", 404)
		}
		minDownPayment, err := strconv.ParseFloat(r.PostFormValue("min-down-payment"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[min-down-payment]", 404)
		}
		loanTerm, err := strconv.Atoi(r.PostFormValue("loan-term"))
		if err != nil {
			http.Error(w, "can`t parse form[loan-term]", 404)
		}

		b := &models.Bank{Name: name, Rate: rate, MaxLoan: maxLoan, MinDownPayment: minDownPayment, LoanTerm: uint(loanTerm)}
		b, err = repo.Create(r.Context(), b)
		if err != nil {
			http.Error(w, "can`t create bank", 404)
		}

		a.logger.Printf("Created: %+v", b)
		a.logger.Println("Try Redirect to: /bank-list/" + name)
		r.Method = http.MethodGet
		http.Redirect(w, r, "/bank-list/"+name, http.StatusMovedPermanently)
	default:
	}
}

func (a *App) BankHandler(w http.ResponseWriter, r *http.Request) {
	if matched, err := regexp.Match(`/bank-list/`, []byte(r.URL.Path)); err != nil && !matched {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	tmpl, err := template.ParseFiles("../web/templates/bank.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	resp := BankResponse{}
	repo := bankrepository.New(a.db, *a.logger)
	a.logger.Printf("%v: %s %s %s ", time.Now(), "BankHandler", r.URL.Path, r.Method)

	switch r.Method {
	case http.MethodGet:
		name := strings.TrimPrefix(r.URL.Path, "/bank-list/")
		b, err := repo.BankByName(r.Context(), name)
		if b.Name == "" || err != nil {
			http.Error(w, "can`t find bank by name", 400)
			return
		}
		a.logger.Println("found bank: ", b.Name)
		resp.Bank = b
		tmpl.Execute(w, resp)
	case http.MethodPost:
		tmpl.Execute(w, resp)
	default:
	}
}

func (a *App) BankCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bank" {
		http.Error(w, "can`t parse form index.html", 404)
		return
	}
	tmpl, err := template.ParseFiles("../web/templates/bank_create.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
	}
	if r.Method != http.MethodGet {
		http.Error(w, "can use only GET method index.html", 404)
		return
	}
	resp := BankResponse{}
	a.logger.Printf("%v: %s %s ", time.Now(), r.URL.Path, r.Method)
	tmpl.Execute(w, resp)
}
