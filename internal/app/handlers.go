package app

import (
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
	MortgageCalcRow         string
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

	tmpl, err := template.ParseFiles("./web/templates/index.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		a.logger.Printf("Error MortgageCalcHandler.ParseFiles: %s", err)
		return
	}

	resp := MortgageCalcResponse{}
	repo := bankrepository.New(a.db, a.logger)
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "can`t parse Form from bank_create.html", 404)
		a.logger.Printf("Error MortgageCalcHandler.ParseForm: %s", err)
		return
	}
	a.logger.Printf("r.Form: %#v", r.Form)
	if r.Form.Has("init-loan") && r.Form.Has("down-payment") && r.Form.Has("name") {
		initLoan, err := strconv.ParseFloat(r.FormValue("init-loan"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[init-loan]", 404)
			a.logger.Printf("Error MortgageCalcHandler.FormValue.init-loan: %s", err)
			return
		}
		downPayment, err := strconv.ParseFloat(r.FormValue("down-payment"), 64)
		if err != nil {
			http.Error(w, "can`t parse form[down-payment]", 404)
			a.logger.Printf("Error MortgageCalcHandler.FormValue.down-payment: %s", err)
			return
		}
		name := r.FormValue("name")

		a.logger.Printf("r.Form: %v,%v,%v", initLoan, downPayment, name)
		b, err := repo.BankByName(r.Context(), name)
		if err != nil {
			http.Error(w, "can`t get bank by name", 404)
			a.logger.Printf("Error MortgageCalcHandler.BankByName: %s", err)
			return
		}
		calcMMP(&resp, b, name, initLoan, downPayment)
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
	tmpl, err := template.ParseFiles("./web/templates/bank_list.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	resp := BankListResponse{}
	repo := bankrepository.New(a.db, a.logger)
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
			http.Error(w, "can`t create bank: Bank name must be unique!", 404)
		}

		a.logger.Printf("Created: %+v", b)
		a.logger.Println("Try Redirect to: /bank-list/" + name)
		http.Redirect(w, r, "/bank-list/"+name, http.StatusSeeOther)
		return
	default:
	}
	err = tmpl.Execute(w, resp)
	if err != nil {
		http.Error(w, "can`t execute template", 404)
		return
	}
}

func (a *App) BankHandler(w http.ResponseWriter, r *http.Request) {
	if matched, err := regexp.Match(`/bank-list/`, []byte(r.URL.Path)); err != nil && !matched {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/bank.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
		return
	}

	resp := BankResponse{}
	repo := bankrepository.New(a.db, a.logger)
	name := strings.TrimPrefix(r.URL.Path, "/bank-list/")
	a.logger.Printf("%v: %s %s %s ", time.Now(), "BankHandler", r.URL.Path, r.Method)

	switch r.Method {
	case http.MethodGet:
		b, err := repo.BankByName(r.Context(), name)
		if b.Name == "" || err != nil {
			http.Error(w, "can`t find bank by name", 400)
			return
		}
		a.logger.Println("found bank: ", b.Name)
		resp.Bank = b
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
		nb, err := repo.Update(r.Context(), b)
		if nb.Name == "" || err != nil {
			http.Error(w, "can`t find bank by name", 400)
			return
		}
		a.logger.Println("bank updated: ", b.Name)

		http.Redirect(w, r, "/bank-list/"+name, http.StatusSeeOther)
		return
	default:
	}
	err = tmpl.Execute(w, resp)
	if err != nil {
		http.Error(w, "can`t execute template", 404)
		return
	}
}

func (a *App) BankCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bank" {
		http.Error(w, "can`t parse form index.html", 404)
		return
	}
	tmpl, err := template.ParseFiles("./web/templates/bank_create.html")
	if err != nil {
		http.Error(w, "can`t parse template index.html", 404)
	}
	if r.Method != http.MethodGet {
		http.Error(w, "can use only GET method index.html", 404)
		return
	}
	resp := BankResponse{}
	a.logger.Printf("%v: %s %s ", time.Now(), r.URL.Path, r.Method)
	err = tmpl.Execute(w, resp)
	if err != nil {
		http.Error(w, "can`t execute template", 404)
		return
	}
}

func (a *App) BankDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bank/delete" {
		http.Error(w, "can`t parse form index.html", 404)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "can use only GET method index.html", 404)
		return
	}

	repo := bankrepository.New(a.db, a.logger)
	a.logger.Printf("%v: %s %s ", time.Now(), r.URL.Path, r.Method)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "can`t parse form index.html", 404)
		return
	}
	a.logger.Printf("%#v", r.PostForm)
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
	err = repo.Delete(r.Context(), b)
	if err != nil {
		http.Error(w, "can`t delete bank", 404)
		return
	}
	http.Redirect(w, r, "/bank-list", http.StatusFound)
}
