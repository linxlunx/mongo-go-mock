package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

type Service struct {
	repository Repository
}

func main() {
	dbURL := "mongodb://127.0.0.1:27017"
	dbName := "ml-gateway-testing"

    db := InitDB(dbURL, dbName)
	s := &Service{
		repository: NewRepository(db),
	}

	newUser := User{
		Username: "user",
		Email: "user@domain.com",
	}

	_id, err := s.repository.Insert(newUser)
	if err != nil {
		panic(err)
	}

	user, err := s.repository.FindOne(_id)
	if err != nil {
		panic(err)
	}
    fmt.Println(user)

	update := User{
		Email: "change@domain.com",
	}

	updated, err := s.repository.Update(user.ID, update)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)
}
