package main

import (
	"RegistrationBotTutorial/src"
)

var telegramBot src.TelegramBot

func main() {
	src.Connection.Init()
	src.Connection.Find(1231)
	//telegramBot.Init()
	//telegramBot.Start()
}
