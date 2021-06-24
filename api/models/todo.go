package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	gorm.Model
	Task string `gorm:"size:300; not null" json:"task"`
	DueDate time.Time `gorm:"not null" json:"due_date"`
	UserID uint `gorm:"not_null" json:"user_id"`
}

func (t *Todo) PrepareTodo() {
	t.Task = strings.TrimSpace(t.Task)
}

func (t *Todo) ValidateTodo() error {
	if t.Task == "" {
		return errors.New("task is required")
	}
	if t.DueDate.IsZero() {
		return errors.New("due date is required")
	}
	return nil
}

func (t *Todo) CreateTodo(db *gorm.DB) (*Todo, error) {
	err := db.Debug().Create(&t).Error
	if err != nil {
		return &Todo{}, err
	}
	return t, nil
}

func (t *Todo) UpdateTodo(id int, db *gorm.DB) (*Todo, error) {
	if err := db.Debug().Table("todos").Where("id = ?", id).Updates(
		Todo {
			Task: t.Task,
		},
	).Error; err != nil {
		return &Todo{}, err
	}
	return t, nil
}

func DeleteTodo(id int, db *gorm.DB) error {
	if err := db.Debug().Table("todos").Where("id = ?", id).Delete(&Todo{}).Error; err != nil {
		return err
	}
	return nil
}

func GetTodoById(id int, db *gorm.DB) (*Todo, error) {
	todo := &Todo{}
	if err := db.Debug().Table("todos").Where("id = ?", id).First(&todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func GetUserTodos(user *User, db *gorm.DB) (*[]Todo, error) {
	todos := []Todo{}
	err := db.Debug().Model(&user).Related(&todos).Error
	if err != nil {
		return &[]Todo{}, err
	}
	return &todos, nil
}