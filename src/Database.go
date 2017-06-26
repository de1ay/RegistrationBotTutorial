package src

import (
	"gopkg.in/mgo.v2"
	"RegistrationBotTutorial/conf"
	"log"
	"gopkg.in/mgo.v2/bson"
)

var Connection DatabaseConnection

type User struct {
	Chat_ID      int64
	Phone_Number string
}

type DatabaseConnection struct {
	Session *mgo.Session  // Соединение с сервером
	DB      *mgo.Database // Соединение с базой данных
}

// Инициализация соединения с БД
func (connection *DatabaseConnection) Init() {
	session, err := mgo.Dial(conf.MONGODB_CONNECTION_URL)
	if err != nil {
		log.Fatal(err)
	}
	connection.Session = session
	db := session.DB(conf.MONGODB_DATABASE_NAME)
	connection.DB = db
}

// Проверка на существование пользователя
func (connection *DatabaseConnection) Find(chatID int64) (bool, error) {
	collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
	count, err := collection.Find(bson.M{"chat_id": chatID}).Count()
	if err != nil || count == 0 {
		return false, err
	} else {
		return true, err
	}
}

// Получение пользователя
func (connection *DatabaseConnection) GetUser(chatID int64) (User, error) {
	var result User
	find, err := connection.Find(chatID)
	if err != nil {
		return result, err
	}
	if find {
		collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
		err = collection.Find(bson.M{"chat_id": chatID}).One(&result)
		return result, err
	} else {
		return result, mgo.ErrNotFound
	}
}

// Создание пользователя
func (connection *DatabaseConnection) CreateUser(user User) error {
	collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
	err := collection.Insert(user)
	return err
}

// Обновление номера мобильного телефона
func (connection *DatabaseConnection) UpdateUser(user User) error {
	collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
	err := collection.Update(bson.M{"chat_id": user.Chat_ID}, &user)
	return err
}
