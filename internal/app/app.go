package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	server http.Server
	db     *sql.DB
	logger *log.Logger
}

func New(config *Config) *App {
	s := initServer(config.Port)
	db := initDatabase(config.DbURL)
	l := initLogger(config.DebugLevel)
	return &App{s, db, l}
}

func (a *App) Run() {
	a.server.Handler = initRoutes(a)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error: ", err)
	}
}
