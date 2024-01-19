package db_test

import (
	"testing"
	"todo-list-service/pkg/db"
	"todo-list-service/pkg/env"
)

type MockMongo struct {
	Close func()
}

func TestMockMongo_HostMemoryDb(t *testing.T) {
	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	t.Run("Successfully host memory mongo", func(t *testing.T) {
		mm := db.MockMongo{}
		_, err := mm.HostMemoryDb(cfg.MongodPath)
		if err != nil {
			t.Error("Failed to connect to memory server")
			t.FailNow()
		}
		defer mm.Close()
	})

}
