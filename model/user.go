package model

import (
	"fmt"
	"time"

	"github.com/sebastian-nunez/golang-search-engine/utils"
	"gorm.io/gorm"
)

type User struct {
	ID        string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `gorm:"default:false" json:"isAdmin"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin(gdb *gorm.DB, email string, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("unable to hash password: %s", err)
	}

	user := User{
		Email:    email,
		Password: hashedPassword,
		IsAdmin:  true,
	}

	if err := gdb.Create(&user).Error; err != nil {
		return fmt.Errorf("unable to create user in the database: %s", err)
	}

	return nil
}

func (u *User) LoginAsAdmin(gdb *gorm.DB, email string, password string) (*User, error) {
	if err := gdb.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if match := utils.ComparePasswords(u.Password, password); !match {
		return nil, fmt.Errorf("invalid password")
	}

	return u, nil
}
