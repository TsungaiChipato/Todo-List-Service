package db

import (
	"context"
	"reflect"
	"testing"
	"time"
	"todo-list-service/pkg/env"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createDb(t *testing.T) (*mongo.Database, func()) {
	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	mm := MockMongo{}
	uri, err := mm.HostMemoryDb(cfg.MongodPath)
	if err != nil {
		t.Error("Failed to launch memory db")
		t.FailNow()
	}

	conn := Connection{}
	err = conn.Connect(uri)
	if err != nil {
		t.Error("Failed to connect to memory server")
		t.Fail()
	}

	return conn.Database, func() {
		conn.Close()
		mm.Close()
	}
}

func createColl(ctx context.Context, t *testing.T) (h TodoItemDbHandler, close func()) {
	db, close := createDb(t)

	h = TodoItemDbHandler{}
	err := h.New(ctx, db)
	if err != nil {
		t.Error("Failed to create the collection")
		t.FailNow()
	}
	return
}

func TestTodoItemDbHandler_New(t *testing.T) {
	t.Parallel()
	db, close := createDb(t)
	defer close()

	t.Run("Successfully create the collection", func(t *testing.T) {
		ctx := context.Background()
		h := TodoItemDbHandler{}

		err := h.New(ctx, db)
		if err != nil {
			t.Errorf("TodoItemDbHandler.New() error = %v, wantErr %v", err, false)
			return
		}
	})
}

func TestTodoItemDbHandler_InsertOne(t *testing.T) {
	t.Parallel()

	t.Run("Successfully inserted one article", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		h, close := createColl(ctx, t)
		defer close()

		item := &TodoItemDb{
			Title:       "Test_Title",
			Description: "Test_Description",
			DueDate:     time.Now().Add(time.Hour).UTC().Truncate(time.Millisecond), // have to truncate, because mongo does not store microseconds
		}
		id, err := h.InsertOne(ctx, item)
		if err != nil {
			t.Errorf("TodoItemDbHandler.InsertOne() error = %v, wantErr %v", err, false)
			return
		}

		createdItem, err := h.FindOneById(ctx, id)
		if err != nil {
			t.Errorf("TodoItemDbHandler.InsertOne() error = %v, wantErr %v", err, false)
			return
		}

		item.Id = createdItem.Id
		if !reflect.DeepEqual(*createdItem, item) {
			t.Errorf("TodoItemDbHandler.InsertOne() = %v, want %v", *createdItem, item)
			return
		}
	})

	t.Run("Successfully inserted one article with predefined labels", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		h, close := createColl(ctx, t)
		defer close()

		item := &TodoItemDb{
			Title:       "Test_Title",
			Description: "Test_Description",
			Labels:      []string{"Test_Label"},
			DueDate:     time.Now().Add(time.Hour).UTC().Truncate(time.Millisecond), // have to truncate, because mongo does not store microseconds
		}
		id, err := h.InsertOne(ctx, item)
		if err != nil {
			t.Errorf("TodoItemDbHandler.InsertOne() error = %v, wantErr %v", err, false)
			return
		}

		createdItem, err := h.FindOneById(ctx, id)
		if err != nil {
			t.Errorf("TodoItemDbHandler.InsertOne() error = %v, wantErr %v", err, false)
			return
		}

		item.Id = createdItem.Id
		if !reflect.DeepEqual(*createdItem, item) {
			t.Errorf("TodoItemDbHandler.InsertOne() = %v, want %v", *createdItem, item)
			return
		}
	})
}

func TestTodoItemDbHandler_FindOneById(t *testing.T) {
	t.Parallel()

	t.Run("Successfully found the image by id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		h, close := createColl(ctx, t)
		defer close()

		item := &TodoItemDb{
			Title:       "Test_Title",
			Description: "Test_Description",
			Labels:      []string{"Test_Label"},
			DueDate:     time.Now().Add(time.Hour).UTC().Truncate(time.Millisecond), // have to truncate, because mongo does not store microseconds
		}
		id, err := h.InsertOne(ctx, item)
		if err != nil {
			t.Errorf("TodoItemDbHandler.FindOneById() error = %v, wantErr %v", err, false)
			return
		}

		createdItem, err := h.FindOneById(ctx, id)
		if err != nil {
			t.Errorf("TodoItemDbHandler.FindOneById() error = %v, wantErr %v", err, false)
			return
		}

		item.Id = createdItem.Id
		if !reflect.DeepEqual(*createdItem, item) {
			t.Errorf("TodoItemDbHandler.FindOneById() = %v, want %v", *createdItem, item)
			return
		}
	})

	t.Run("Successfully found nothing with non-existing article", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		h, close := createColl(ctx, t)
		defer close()

		found, err := h.FindOneById(ctx, primitive.NewObjectID())
		if err != nil {
			t.Errorf("TodoItemDbHandler.FindOneById() error = %v, wantErr %v", err, false)
			return
		}

		if found != nil {
			t.Errorf("TodoItemDbHandler.FindOneById() = %v, want %v", *found, nil)
			return
		}
	})
}

// TODO: other db tests
