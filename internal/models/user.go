package models

import (
	"golang.org/x/crypto/bcrypt"
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

func (u *User) ValidatePassword(input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input)); err != nil {
		return false
	}

	return true
}
