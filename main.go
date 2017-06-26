package main

import (
	"RegistrationBotTutorial/src"
)

var telegramBot src.TelegramBot

func main() {
	src.Connection.Init()
	telegramBot.Init()
	telegramBot.Start()
}
