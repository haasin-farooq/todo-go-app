package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/haasin-farooq/todo-go-app/api/models"
	"github.com/haasin-farooq/todo-go-app/api/responses"
	"github.com/haasin-farooq/todo-go-app/utils"
)

func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "success",
		"message": "Registration successful",
	}

	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	u, _ := user.GetUser(a.DB)
	if u != nil {
		res["status"] = "failed"
		res["message"] = "User already exists, please login"
		responses.JSON(w, http.StatusBadRequest, res)
		return
	}

	user.PrepareUser()

	err = user.ValidateUser("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userCreated, err := user.CreateUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res["user"] = userCreated
	responses.JSON(w, http.StatusOK, res)
}

func (a *App) UserLogin(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"status": "success",
		"message": "Login successful",
	}

	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.PrepareUser()

	err = user.ValidateUser("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	u, err := user.GetUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if u == nil {
		res["status"] = "failed"
		res["message"] = "User not found, please register"
		responses.JSON(w, http.StatusBadRequest, res)
		return
	}

	err = models.CheckPasswordHash(user.Password, u.Password)
	if err != nil {
		res["status"] = "failed"
		res["message"] = "Login failed, please try again"
		responses.JSON(w, http.StatusBadRequest, res)
		return
	}

	token, err := utils.EncodeAuthToken(u.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res["token"] = token
	responses.JSON(w, http.StatusOK, res)
}