package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Repository interface {
	Insert(userData User) (primitive.ObjectID, error)
	Update(_id primitive.ObjectID, updateData User) (*User, error)
	FindOne(_id primitive.ObjectID) (*User, error)
	FindAll() ([]User, error)
}

type CollRepository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *CollRepository {
	coll := db.Collection("users")
	return &CollRepository{
		coll: coll,
	}
}

func (r *CollRepository) Insert(userData User) (primitive.ObjectID, error) {
	result, err := r.coll.InsertOne(context.TODO(), userData)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *CollRepository) Update(_id primitive.ObjectID, updateData User) (*User, error) {
	err := r.coll.FindOneAndUpdate(
		context.TODO(),
		bson.D{{"_id", _id}},
		bson.D{{"$set", updateData}},
		options.FindOneAndUpdate().SetReturnDocument(1),
	).Decode(&updateData)

	if err != nil {
		return nil, err
	}
	return &updateData, nil
}

func (r *CollRepository) FindOne(_id primitive.ObjectID) (*User, error) {
	var user User
	err := r.coll.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *CollRepository) FindAll() ([]User, error) {
	var users []User

	cur, err := r.coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return []User{}, err
	}

	for cur.Next(context.TODO()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}