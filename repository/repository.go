package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	database *mongo.Database
}

func GetNew(database *mongo.Database) *Repository {
	return &Repository{
		database,
	}
}
