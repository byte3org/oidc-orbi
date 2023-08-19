package models

import (
	"github.com/byte3org/oidc-orbi/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ModelBase
	Username      string
	Password      string
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	Phone         string
	PhoneVerified bool
	IsAdmin       bool
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hash, err := utils.MakePassword(u.Password)
		if err != nil {
			return err
		}

		tx.Statement.SetColumn("Username", hash)
	}

	return
}

func (u *User) ValidatePassword(input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input)); err != nil {
		return false
	}

	return true
}
