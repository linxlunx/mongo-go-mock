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

func TestFindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T){
		db := mt.Client.Database("test_db")
		s := &Service{
			repository: NewRepository(db),
		}

		users := []User{
			{
				ID: primitive.NewObjectID(),
				Email: "first@first.com",
				Username: "first",
			},
			{
				ID: primitive.NewObjectID(),
				Email: "second@second.com",
				Username: "",
			},
		}

		first := mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, bson.D{
			{"_id", users[0].ID},
			{"email", users[0].Email},
			{"username", users[0].Username},
		})

		second := mtest.CreateCursorResponse(1, "test_db.users", mtest.NextBatch, bson.D{
			{"_id", users[1].ID},
			{"email", users[1].Email},
			{"username", users[1].Username},
		})

		killCursors := mtest.CreateCursorResponse(0,"test_db.users", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		userList, err := s.repository.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, users, userList)
	})
}

func TestInsert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T){
		db := mt.Client.Database("test_db")
		s := &Service{
			repository: NewRepository(db),
		}
		_id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		insertedID, err := s.repository.Insert(User{
			ID: _id,
			Email: "user@user.com",
			Username: "user",
		})

		assert.Nil(t, err)
		assert.Equal(t, _id, insertedID)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T){
		db := mt.Client.Database("test_db")
		s := &Service{
			repository: NewRepository(db),
		}

		_id := primitive.NewObjectID()
		updateData := User{
			Email: "second@email.com",
			Username: "second",
		}

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"_id", _id},
				{"email", updateData.Email},
				{"username", updateData.Username},
			}},
		})

		updated, err := s.repository.Update(_id, updateData)
		updateData.ID = _id
		assert.Nil(t, err)
		assert.Equal(t, updated, &updateData)
	})
}

