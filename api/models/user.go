package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index" json:"email"`
	FirstName string `gorm:"size:100;not null" json:"first_name"`
	LastName string `gorm:"size:100;not null" json:"last_name"`
	Password string `gorm:"size:100;not null" json:"password"`
}

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashPassword), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

// BeforeSave Hook is called automatically
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) PrepareUser() {
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
}

func (u *User) ValidateUser(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil
	default:
		if u.FirstName == "" {
			return errors.New("first_name is required")
		}
		if u.LastName == "" {
			return errors.New("last_name is required")
		}
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) GetUser(db *gorm.DB) (*User, error) {
	user := &User{}
	if err := db.Debug().Table("users").Where("email = ?", u.Email).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}