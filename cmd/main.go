package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var srv http.Server

	http.HandleFunc("/", SpinInDocker)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error: ", err)
	}

}

func SpinInDocker(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\nSpin in Docker container and Run throw Docker-compose!!")
}
