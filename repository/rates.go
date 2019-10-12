package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LordotU/my-savings-telegram-bot/types"
)

const RATES_COLLECTION = "rates"

func (repository *Repository) FindRate(currency string) (*types.Rate, error) {
	rate := types.Rate{}
	err := repository.database.Collection(RATES_COLLECTION).FindOne(
		context.TODO(),
		bson.M{"_id": currency},
	).Decode(&rate)
	return &rate, err
}

func (repository *Repository) FindRates(currencies []string) ([]*types.Rate, error) {
	filter := bson.M{}
	if len(currencies) > 0 {
		filter = bson.M{"_id": bson.M{"$in": currencies}}
	}
	cursor, err := repository.database.Collection(RATES_COLLECTION).Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{"_id", 1}}),
	)
	if err != nil {
		return nil, err
	}

	var rates []*types.Rate

	for cursor.Next(context.TODO()) {
		var rate types.Rate
		if err := cursor.Decode(&rate); err != nil {
			return nil, err
		}
		rates = append(rates, &rate)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rates, nil
}

func (repository *Repository) UpdateRates(exchangeRates []*types.Rate) (err error) {
	for _, rate := range exchangeRates {
		existingRate, err := repository.FindRate(rate.Currency)
		if err != nil && err.Error() != "mongo: no documents in result" {
			return err
		}

		if existingRate != nil && existingRate.Currency != "" {
			_, err = repository.database.Collection(RATES_COLLECTION).UpdateOne(
				context.TODO(),
				bson.M{"_id": rate.Currency},
				bson.M{"$set": bson.M{"rate": rate.Rate, "updated_at": time.Now()}},
			)
		} else {
			_, err = repository.database.Collection(RATES_COLLECTION).InsertOne(
				context.TODO(),
				rate,
			)
		}

		if err != nil {
			return err
		}
	}

	return
}
