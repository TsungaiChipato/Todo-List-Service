package router

import (
	"todo-list-service/pkg/controller"
	"todo-list-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func AttachTodoItemRoutes(engine *gin.Engine, ctrl *controller.TodoItemController) {
	engine.GET("/todo", ctrl.FindAll)
	engine.GET("/todo/:id", middleware.IdParam(), ctrl.FindOneById)
	engine.GET("/todo/label/:label", middleware.LabelParam(), ctrl.FindByLabel)

	engine.POST("/todo", ctrl.Create)

	engine.PUT("/todo/:id", middleware.IdParam(), ctrl.UpdateByID)

	engine.DELETE("/todo/:id", middleware.IdParam(), ctrl.DeleteOneById)
}
