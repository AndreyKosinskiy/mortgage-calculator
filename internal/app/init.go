package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

func initServer(port string) http.Server {
	s := http.Server{
		Addr: port,
	}
	return s
}

func initRoutes(a *App) *http.ServeMux {
	route := http.DefaultServeMux

	fs := http.FileServer(http.Dir("../web/static/"))
	route.Handle("/static/", http.StripPrefix("/static/", fs))

	route.HandleFunc("/", a.MortgageCalcHandler)
	route.HandleFunc("/bank", a.BankCreateHandler)
	route.HandleFunc("/bank-list", a.BankListHandler)
	route.HandleFunc("/bank-list/", a.BankHandler)
	return route
}

func initDatabase(settings string) *sql.DB {
	return nil
}

func initLogger(debugLevel int) *log.Logger {
	return log.New(os.Stdout, "Debug: ", 0)
}
