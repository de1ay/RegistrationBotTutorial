package src

import (
	"gopkg.in/mgo.v2"
	"RegistrationBotTutorial/conf"
	"log"
	"gopkg.in/mgo.v2/bson"
)

var Connection DatabaseConnection

type User struct {
	ChatID      int64 `json:"chat_id"`
	PhoneNumber string `json:"phone_number"`
}

type DatabaseConnection struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (connection *DatabaseConnection) Init() {
	session, err := mgo.Dial(conf.MONGODB_CONNECTION_URL)
	if err != nil {
		log.Fatal(err)
	}
	connection.Session = session
	db := session.DB(conf.MONGODB_DATABASE_NAME)
	connection.DB = db
}

func (connection *DatabaseConnection) Find(chatID int64) (bool, error) {
	collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
	count, err := collection.Find(bson.M{"chat_id": chatID}).Count()
	if err != nil || count == 0 {
		return false, err
	} else {
		return true, err
	}
}

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

func (connection *DatabaseConnection) CreateUser(user User) error {
	collection := connection.DB.C(conf.MONGODB_COLLECTION_USERS)
	err := collection.Insert(user)
	return err
}
