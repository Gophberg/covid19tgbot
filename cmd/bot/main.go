package main

import (
	"covid19tgbot/internal/app/bot"
)

func main() {
	config := bot.NewConfig()
	bot.DBInit(config)
	bot.RunBot(config)
}
