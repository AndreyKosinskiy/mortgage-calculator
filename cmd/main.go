package main

import (
	"github.com/AndreyKosinskiy/mortgage-calculator/internal/app"
)

func main() {
	config := app.NewConfig()
	app := app.New(config)
	app.Run()
}
