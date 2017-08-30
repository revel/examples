package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model
	Name           string `gorm:"size:255"`
	Email          string `gorm:"type:varchar(100);unique_index"`
	HashedPassword []byte
	Active         bool
	FileName       string `gorm:"size:255"`
}

// SetNewPassword set a new hashsed password to user
func (user *User) SetNewPassword(passwordString string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	user.HashedPassword = bcryptPassword
}
