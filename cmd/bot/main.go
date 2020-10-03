package main

import (
	"covid19rus/internal/bot"
)

func main() {
	config := bot.NewConfig()
	bot.DBInit(config)
	bot.RunBot(config)
}
