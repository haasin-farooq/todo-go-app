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

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/haasin-farooq/todo-go-app/api/controllers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := controllers.App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	app.RunServer()
}