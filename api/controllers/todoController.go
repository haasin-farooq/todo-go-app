package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/haasin-farooq/todo-go-app/api/models"
	"github.com/haasin-farooq/todo-go-app/api/responses"
)

func (a *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "success",
		"message": "Todo created successfully",
	}

	u := r.Context().Value("userID").(float64)
	userID := uint(u)

	todo := &models.Todo{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &todo)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo.PrepareTodo()

	err = todo.ValidateTodo()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo.UserID = userID

	todoCreated, err := todo.CreateTodo(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res["todo"] = todoCreated
	responses.JSON(w, http.StatusCreated, res)
}

func (a *App) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "success",
		"message": "Todo updated successfully",
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	u := r.Context().Value("userID").(float64)
	userID := uint(u)

	t, err := models.GetTodoById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if t == nil {
		res["status"] = "failed"
		res["message"] = "Todo not found"
		responses.JSON(w, http.StatusBadRequest, res)
		return
	}

	if t.UserID != userID {
		res["status"] = "failed"
		res["message"] = "Unauthorized action"
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo := &models.Todo{}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	todo.PrepareTodo()

	err = todo.ValidateTodo()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	updatedTodo, err := todo.UpdateTodo(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	updatedTodo.ID = uint(id)
	updatedTodo.UserID = userID

	res["todo"] = updatedTodo
	responses.JSON(w, http.StatusOK, res)
}

func (a *App) GetUserTodos(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "failed",
		"message": "Unauthorized action",
	}

	u := r.Context().Value("userID").(float64)
	userID := uint(u)

	user, err := models.GetUserById(userID, a.DB)
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