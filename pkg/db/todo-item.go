package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoItemDbHandler struct {
	coll *mongo.Collection
}

type TodoItemDbHandlerInterface interface {
	New(context.Context, *mongo.Database) error
	InsertOne(context.Context, *TodoItemDb) (primitive.ObjectID, error)
	FindOneById(context.Context, primitive.ObjectID) (*TodoItemDb, error)
	FindAll(context.Context) (*mongo.Cursor, error)
	FindByLabel(context.Context, string) (*mongo.Cursor, error)
	AddLabel(context.Context, primitive.ObjectID, string) error
	RemoveLabel(context.Context, primitive.ObjectID, string) error
	DeleteOneById(context.Context, primitive.ObjectID) error
	UpdateOneById(context.Context, primitive.ObjectID, *TodoItemDb) error
	ConsumeCursor(*mongo.Cursor, int) (*[]TodoItemDb, error)
}

type TodoItemDb struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"string" json:"string"`
	DueDate     time.Time          `bson:"dueDate" json:"dueDate"`
	Labels      []string           `bson:"labels,omitempty" json:"labels,omitempty"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed,omitempty" json:"completed,omitempty"`
}

func (h *TodoItemDbHandler) New(context context.Context, database *mongo.Database) error {
	h.coll = database.Collection("articles")

	// creates an index on the labels array
	_, err := h.coll.Indexes().CreateOne(context, mongo.IndexModel{
		Keys: bson.D{{Key: "labels", Value: 1}},
	})
	return err
}

func (h *TodoItemDbHandler) InsertOne(context context.Context, new *TodoItemDb) (primitive.ObjectID, error) {
	result, err := h.coll.InsertOne(context, new)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (h *TodoItemDbHandler) FindOneById(context context.Context, id primitive.ObjectID) (*TodoItemDb, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var article TodoItemDb
	err := h.coll.FindOne(context, filter).Decode(&article)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &article, nil
}

func (h *TodoItemDbHandler) FindAll(context context.Context) (*mongo.Cursor, error) {
	cur, err := h.coll.Find(context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	return cur, nil
}

func (h *TodoItemDbHandler) FindByLabel(context context.Context, label string) (*mongo.Cursor, error) {
	filter := bson.M{"labels": label}
	cur, err := h.coll.Find(context, filter)
	if err != nil {
		return nil, err
	}

	return cur, nil
}

func (h *TodoItemDbHandler) AddLabel(context context.Context, id primitive.ObjectID, label string) error {
	update := bson.M{"$addToSet": bson.M{"labels": label}} // should not have duplicate labels
	_, err := h.coll.UpdateByID(context, id, update)
	return err
}

func (h *TodoItemDbHandler) RemoveLabel(context context.Context, id primitive.ObjectID, label string) error {
	update := bson.M{"$pull": bson.M{"labels": label}}
	_, err := h.coll.UpdateByID(context, id, update)
	return err
}

func (h *TodoItemDbHandler) DeleteOneById(context context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := h.coll.DeleteOne(context, filter)
	return err
}

func (h *TodoItemDbHandler) UpdateOneById(context context.Context, id primitive.ObjectID, update *TodoItemDb) error {
	_, err := h.coll.UpdateByID(context, id, bson.M{"$set": update})
	return err
}

func (h *TodoItemDbHandler) ConsumeCursor(cur *mongo.Cursor, max int) (*[]TodoItemDb, error) {
	results := []TodoItemDb{}

	i := 0
	for cur.Next(context.TODO()) {
		var elem TodoItemDb
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
		i++

		if i == max {
			break
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())
	return &results, nil
}
