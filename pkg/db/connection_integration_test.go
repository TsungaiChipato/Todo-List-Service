package db_test

import (
	"testing"
	"todo-list-service/pkg/db"
	"todo-list-service/pkg/env"
)

func TestConnection_Connect(t *testing.T) {
	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	t.Run("Successfully connect to memory mongo", func(t *testing.T) {
		mm := db.MockMongo{}
		uri, err := mm.HostMemoryDb(cfg.MongodPath)
		if err != nil {
			t.Error("Failed to launch memory db")
			t.FailNow()
		}

		defer mm.Close()

		conn := db.Connection{}
		err = conn.Connect(uri)
		defer conn.Close()

		if err != nil {
			t.Error("Failed to connect to memory server")
			t.Fail()
		}
	})

}
