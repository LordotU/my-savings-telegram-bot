package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LordotU/my-savings-telegram-bot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const USERS_COLLECTION = "users"

func (repository *Repository) CreateUser(tgUser tgbotapi.User, baseCurrency string) (user *mongo.InsertOneResult, err error) {
	user, err = repository.database.Collection(USERS_COLLECTION).InsertOne(
		context.TODO(),
		types.GetNewUser(tgUser, baseCurrency),
	)
	return
}

func (repository *Repository) FindUser(telegram_id int) (*types.User, error) {
	user := types.User{}
	err := repository.database.Collection(USERS_COLLECTION).FindOne(
		context.TODO(),
		bson.M{"telegram_id": telegram_id},
	).Decode(&user)
	return &user, err
}

func (repository *Repository) UpdateUser(telegram_id int, data map[string]interface{}) (err error) {
	data["updated_at"] = time.Now()

	_, err = repository.database.Collection(USERS_COLLECTION).UpdateOne(
		context.TODO(),
		bson.M{"telegram_id": telegram_id},
		bson.M{"$set": data},
	)

	return
}
