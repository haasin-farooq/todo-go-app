package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TempTodo struct {
	Task string `gorm:"size:300; not null" json:"task"`
	DueDate time.Time `gorm:"not null" json:"due_date"`
	Email string `gorm:"type:varchar(100); not null" json:"email"`
}

func (tt *TempTodo) CreateTempTodo(db *gorm.DB) (*TempTodo, error) {
	err := db.Debug().Create(&tt).Error
	if err != nil {
		return &TempTodo{}, err
	}
	return tt, nil
}

func GetTempTodosByEmail(email string, db *gorm.DB) ([]TempTodo, error) {
	tt := []TempTodo{}
	if err := db.Debug().Table("temp_todos").Where("email = ?", email).Find(&tt).Error; err != nil {
		return []TempTodo{}, err
	}
	return tt, nil
}

func DeleteTempTodos(email string, db *gorm.DB) error {
	if err := db.Debug().Table("temp_todos").Where("email = ?", email).Delete(&TempTodo{}).Error; err != nil {
		return err
	}
	return nil
}