package src

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"RegistrationBotTutorial/conf"
)

type TelegramBot struct {
	API                   *tgbotapi.BotAPI        // API телеграмма
	Updates               tgbotapi.UpdatesChannel // Канал обновлений
	ActiveContactRequests []int64                 // ID чатов, от которых мы ожидаем номер
}

// Инициализация бота
func (telegramBot *TelegramBot) Init() {
	botAPI, err := tgbotapi.NewBotAPI(conf.TELEGRAM_BOT_API_KEY)  // Инициализация API
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.API = botAPI
	botUpdate := tgbotapi.NewUpdate(conf.TELEGRAM_BOT_UPDATE_OFFSET)  // Инициализация канала обновлений
	botUpdate.Timeout = conf.TELEGRAM_BOT_UPDATE_TIMEOUT
	botUpdates, err := telegramBot.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.Updates = botUpdates
}

// Основной цикл бота
func (telegramBot *TelegramBot) Start() {
	for update := range telegramBot.Updates {
		if update.Message != nil {
			// Если сообщение есть и его длина больше 0 -> начинаем обработку
			telegramBot.analyzeUpdate(update)
		}
	}
}

func (telegramBot *TelegramBot) addContactRequestID(chatID int64) {
	telegramBot.ActiveContactRequests = append(telegramBot.ActiveContactRequests, chatID)
}

func (telegramBot *TelegramBot) findContactRequestID(chatID int64) bool {
	for _, v := range telegramBot.ActiveContactRequests {
		if v == chatID {
			return true
		}
	}
	return false
}

func (telegramBot *TelegramBot) deleteContactRequestID(chatID int64) {
	for i, v := range telegramBot.ActiveContactRequests {
		if v == chatID {
			copy(telegramBot.ActiveContactRequests[i:], telegramBot.ActiveContactRequests[i + 1:])
			telegramBot.ActiveContactRequests[len(telegramBot.ActiveContactRequests) - 1] = 0
			telegramBot.ActiveContactRequests = telegramBot.ActiveContactRequests[:len(telegramBot.ActiveContactRequests) - 1]
		}
	}
}

// Начало обработки сообщения
func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if telegramBot.findUser(chatID) {
		telegramBot.analyzeUser(update)
	} else {
		telegramBot.createUser(User{chatID, ""})
		telegramBot.requestContact(chatID)
	}
}

// Есть ли пользователь в БД?
func (telegramBot *TelegramBot) findUser(chatID int64) bool {
	find, err := Connection.Find(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
	}
	return find
}

// Создать нового пользователя
func (telegramBot *TelegramBot) createUser(user User) {
	err := Connection.CreateUser(user)
	if err != nil {
		msg := tgbotapi.NewMessage(user.Chat_ID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
	}
}

func (telegramBot *TelegramBot) analyzeUser(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	user, err := Connection.GetUser(chatID)  // Вытаскиваем данные из БД для проверки номера
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
		return
	}
	if len(user.Phone_Number) > 0 {
		msg := tgbotapi.NewMessage(chatID, "Ваш номер: " + user.Phone_Number)  // Если номер у нас уже есть, то пишем его
		telegramBot.API.Send(msg)
		return
	} else {
		// Если номера нет, то проверяем ждём ли мы контакт от этого ChatID
		if telegramBot.findContactRequestID(chatID) {
			telegramBot.checkRequestContactReply(update)  // Если да -> проверяем
			return
		} else {
			telegramBot.requestContact(chatID)  // Если нет -> запрашиваем его
			return
		}
	}
}

// Запросить номер телефона
func (telegramBot *TelegramBot) requestContact(chatID int64) {
	requestContactMessage := tgbotapi.NewMessage(chatID, "Согласны ли вы предоставить ваш номер телефона для регистрации в системе?")
	acceptButton := tgbotapi.NewKeyboardButtonContact("Да")
	declineButton := tgbotapi.NewKeyboardButton("Нет")
	requestContactReplyKeyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{acceptButton, declineButton})
	requestContactMessage.ReplyMarkup = requestContactReplyKeyboard
	telegramBot.API.Send(requestContactMessage)
	telegramBot.addContactRequestID(chatID)
}

// Проверка принятого контакта
func (telegramBot *TelegramBot) checkRequestContactReply(update tgbotapi.Update) {
	if update.Message.Contact != nil {  // Проверяем, содержит ли сообщение контакт
		if update.Message.Contact.UserID == update.Message.From.ID {  // Проверяем действительно ли это контакт отправителя
			telegramBot.updateUser(User{update.Message.Chat.ID, update.Message.Contact.PhoneNumber}, update.Message.Chat.ID)  // Обновляем номер
			telegramBot.deleteContactRequestID(update.Message.Chat.ID)  // Удаляем ChatID из списка ожидания
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)  // Убираем клавиатуру
			telegramBot.API.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер телефона, который вы предоставили, принадлежит не вам!")
			telegramBot.API.Send(msg)
			telegramBot.requestContact(update.Message.Chat.ID)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Если вы не предоставите ваш номер телефона, вы не сможете пользоваться системой!")
		telegramBot.API.Send(msg)
		telegramBot.requestContact(update.Message.Chat.ID)
	}
}

// Обновление номера мобильного телефона пользователя
func (telegramBot *TelegramBot) updateUser(user User, chatID int64) {
	err := Connection.UpdateUser(user)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
		return
	}
}
