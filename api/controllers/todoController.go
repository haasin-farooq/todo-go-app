package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/haasin-farooq/todo-go-app/api/models"
	"github.com/haasin-farooq/todo-go-app/api/responses"
)

func (a *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "success",
		"message": "Todo created successfully",
	}

	userID := r.Context().Value("userID").(float64)

	todo := &models.Todo{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &todo)
	if err != nil {
		fmt.Println(err)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo.PrepareTodo()

	err = todo.ValidateTodo()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo.UserID = int(userID)
	// todo.DueDate, err = time.Parse("_2-_01-2006 03:04 PM", todo.DueDate.String())
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	todoCreated, err := todo.CreateTodo(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res["todo"] = todoCreated
	responses.JSON(w, http.StatusCreated, res)
}

func (a *App) GetUserTodos(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "failed",
		"message": "Unauthorized action, please login",
	}

	userID := r.Context().Value("userID").(float64)

	user, err := models.GetUserById(int(userID), a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if user == nil {
		responses.JSON(w, http.StatusUnauthorized, res)
		return
	}

	todos, err := models.GetUserTodos(user, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, todos)
}