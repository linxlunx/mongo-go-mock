package main

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T){
		db := mt.Client.Database("test_db")
		s := &Service{
			repository: NewRepository(db),
		}

		expected := User{
			ID: primitive.NewObjectID(),
			Email: "admin@admin.com",
			Username: "admin",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, bson.D{
			{"_id", expected.ID},
			{"username", expected.Username},
			{"email", expected.Email},
		}))

		result, err := s.repository.FindOne(expected.ID)
		assert.Nil(t, err)
		assert.Equal(t, &expected, result)
	})
}
