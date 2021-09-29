package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func InitDB(uri string, dbname string) *mongo.Database {
	option := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(option)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())

	db = client.Database(dbname)
	return db
}