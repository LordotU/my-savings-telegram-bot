package main

import "github.com/LordotU/my-savings-telegram-bot/app"

func main() {
	if application, err := app.New(); err != nil {
		panic(err)
	} else {
		application.Run()
	}

}
