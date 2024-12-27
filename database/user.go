package database

import (
	"fmt"
	"time"

	"github.com/sebastian-nunez/golang-search-engine/utils"
)

type User struct {
	ID        string     `gorm:"type:uuid,default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `gorm:"default:false" json:"isAdmin"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return fmt.Errorf("unable to hash password: %s", err)
	}

	// TODO: For testing, always create a set of basic credentials. Remove for a production app.
	user := User{
		Email:    "jdoe@google.com",
		Password: hashedPassword,
		IsAdmin:  true,
	}

	if err := DBConn.Create(&user).Error; err != nil {
		return fmt.Errorf("unable to create user in the database: %s", err)
	}

	return nil
}

func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if match := utils.ComparePasswords(u.Password, password); !match {
		return nil, fmt.Errorf("invalid password")
	}

	return u, nil
}
