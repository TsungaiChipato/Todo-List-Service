package main

import (
	"context"
	"fmt"
	"todo-list-service/pkg/controller"
	"todo-list-service/pkg/db"
	"todo-list-service/pkg/env"
	"todo-list-service/pkg/middleware"
	"todo-list-service/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	var uri string
	if cfg.UseMemoryMongo {
		mm := db.MockMongo{}
		uri, err = mm.HostMemoryDb(cfg.MongodPath)
		if err != nil {
			panic(err)
		}
		defer mm.Close()
	} else {
		uri = cfg.MongoURL
	}

	conn := db.Connection{}
	err = conn.Connect(uri)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	engine := gin.Default()
	engine.Use(middleware.ErrorHandler())

	dbHandler := &db.TodoItemDbHandler{}
	err = dbHandler.New(context.TODO(), conn.Database)
	if err != nil {
		panic(err)
	}

	articleController := &controller.TodoItemController{
		TodoItemDbHandler:  dbHandler,
		MaxReturnArraySize: cfg.MaxReturnArraySize,
	}

	router.AttachTodoItemRoutes(engine, articleController)

	engine.Run(fmt.Sprintf(":%v", cfg.Port))
}
