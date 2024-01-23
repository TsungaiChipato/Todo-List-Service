package controller

import (
	"fmt"
	"net/http"
	"time"
	"todo-list-service/pkg/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoItemController struct {
	TodoItemDbHandler  db.TodoItemDbHandlerInterface
	MaxReturnArraySize int
}

type NewTodoItemBody struct {
	Title       string    `json:"title" binding:"required"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
	Labels      []string  `json:"labels,omitempty"`
	Description string    `json:"description,omitempty" field:"''"`
	Completed   bool      `json:"completed,omitempty" field:"false"`
}

func (con *TodoItemController) FindOneById(c *gin.Context) {
	idString := c.GetString("id")
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to decode id"))
		return
	}

	item, err := con.TodoItemDbHandler.FindOneById(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if item == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (con *TodoItemController) FindAll(c *gin.Context) {
	cur, err := con.TodoItemDbHandler.FindAll(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	items, err := con.TodoItemDbHandler.ConsumeCursor(cur, con.MaxReturnArraySize)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (con *TodoItemController) FindByLabel(c *gin.Context) {
	labelString := c.GetString("label")
	cur, err := con.TodoItemDbHandler.FindByLabel(c, labelString)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	items, err := con.TodoItemDbHandler.ConsumeCursor(cur, con.MaxReturnArraySize)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (con *TodoItemController) DeleteOneById(c *gin.Context) {
	idString := c.GetString("id")
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to decode id"))
		return
	}

	err = con.TodoItemDbHandler.DeleteOneById(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// FIXME: check if the todo item exists
func (con *TodoItemController) UpdateByID(c *gin.Context) {
	idString := c.GetString("id")
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to decode id"))
		return
	}

	todoItem := &NewTodoItemBody{}
	if err := c.ShouldBindJSON(todoItem); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = con.TodoItemDbHandler.UpdateOneById(c, id, &db.TodoItemDb{
		Title:       todoItem.Title,
		DueDate:     todoItem.DueDate,
		Labels:      todoItem.Labels,
		Description: todoItem.Description,
		Completed:   todoItem.Completed,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	item, err := con.TodoItemDbHandler.FindOneById(c, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create is a method of TodoItemController that handles the creation of a new TodoItem.
// It first binds the incoming JSON body to a NewTodoItemBody struct. If there's an error in binding,
// it responds with a 400 status code and the error message.
// If the binding is successful, it attempts to insert a new TodoItem into the database with the data from the struct.
// If there's an error in inserting the data, it responds with a 500 status code.
// If the data is successfully inserted, it responds with a 201 status code and the ID of the newly created TodoItem.
func (con *TodoItemController) Create(c *gin.Context) {
	todoItem := &NewTodoItemBody{}
	if err := c.ShouldBindJSON(todoItem); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := con.TodoItemDbHandler.InsertOne(c, &db.TodoItemDb{
		Title:       todoItem.Title,
		DueDate:     todoItem.DueDate,
		Labels:      todoItem.Labels,
		Description: todoItem.Description,
		Completed:   todoItem.Completed,
	})

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}
