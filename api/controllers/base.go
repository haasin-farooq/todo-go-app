/*
Build a simple Application which stores ToDo tasks with a due date for each user.
The App should have jwt based authentication.
You can use either use MySQL or PostgreSQL (recommended) as the database.
The App should expose following endpoints for CRUD and auth operations:
- Register a new user
- User login
- Add a Task
- Edit a Task
- Delete a Task
- Get All Tasks for a user
- Assign a Internal / External user a task by email address.
  If the user doesn’t exist send them an email to sign up.
  Once they signup that note should be assigned to them automatically.
*/

package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/haasin-farooq/todo-go-app/api/middlewares"
	"github.com/haasin-farooq/todo-go-app/api/models"
	"github.com/haasin-farooq/todo-go-app/api/responses"
)

type App struct {
	Router *mux.Router
	DB *gorm.DB
}

func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	a.DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("\nCannot connect to the database %s\n", DbName)
		log.Fatal("Error: ", err)
	} else {
		fmt.Printf("Connected to the database %s\n", DbName)
	}

	a.DB.Debug().AutoMigrate(&models.User{})

	a.Router = mux.NewRouter().StrictSlash(true)

	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
	a.Router.HandleFunc("/login", a.UserLogin).Methods("POST")
}

func (a *App) RunServer() {
	log.Printf("Server starting on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Todo App!")
}