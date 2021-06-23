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

	a.DB.Debug().AutoMigrate(&models.User{}, &models.Todo{}, &models.TempTodo{})

	a.DB.Debug().Model(&models.Todo{}).AddForeignKey("user_id", "users(id)", "CASCADE", "NO ACTION")

	a.Router = mux.NewRouter().StrictSlash(true)

	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
	a.Router.HandleFunc("/login", a.UserLogin).Methods("POST")
	a.Router.HandleFunc("/todos/{email}", a.CreateTodoByEmail).Methods("POST")

	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middlewares.AuthJwtVerify)

	s.HandleFunc("/todos", a.CreateTodo).Methods("POST")
	s.HandleFunc("/todos", a.GetUserTodos).Methods("GET")
	s.HandleFunc("/todos/{id:[0-9]+}", a.UpdateTodo).Methods("PATCH")
	s.HandleFunc("/todos/{id:[0-9]+}", a.DeleteTodo).Methods("DELETE")
}

func (a *App) RunServer() {
	log.Printf("Server starting on port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Todo App!")
}