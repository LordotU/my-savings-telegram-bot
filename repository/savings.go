package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LordotU/my-savings-telegram-bot/types"
)

const SAVINGS_COLLECTION = "savings"

func (repository *Repository) CreateSaving(
	telegramID int,
	amount float64,
	currency string,
) (err error) {
	_, err = repository.database.Collection(SAVINGS_COLLECTION).InsertOne(
		context.TODO(),
		types.GetNewSaving(telegramID, amount, currency),
	)
	return
}

func (repository *Repository) FindSavings(telegramID int) ([]*types.Saving, error) {
	cursor, err := repository.database.Collection(SAVINGS_COLLECTION).Find(
		context.TODO(),
		bson.M{"telegram_id": telegramID},
		options.Find().SetSort(bson.D{{"updated_at", -1}}),
	)
	if err != nil {
		return nil, err
	}

	var savings []*types.Saving

	for cursor.Next(context.TODO()) {
		var saving types.Saving
		if err := cursor.Decode(&saving); err != nil {
			return nil, err
		}
		savings = append(savings, &saving)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return savings, nil
}

func (repository *Repository) DeleteSavings(savingsIDs []string) (err error) {
	if len(savingsIDs) == 0 {
		return errors.New("savingsIDs for deletion cannot be an empty array!")
	}

	objectsIDs := make([]primitive.ObjectID, len(savingsIDs))
	for i, savingID := range savingsIDs {
		objectsIDs[i], err = primitive.ObjectIDFromHex(savingID)
		if err != nil {
			return err
		}
	}

	filter := bson.M{"_id": bson.M{"$in": objectsIDs}}

	_, err = repository.database.Collection(SAVINGS_COLLECTION).DeleteMany(
		context.TODO(),
		filter,
	)

	return
}
