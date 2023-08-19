package storage

import (
	"crypto/rsa"
	"log"

	"github.com/byte3org/oidc-orbi/internal/models"
	"github.com/byte3org/oidc-orbi/internal/repository"
)

type Service struct {
	keys map[string]*rsa.PublicKey
}

type UserStore interface {
	GetUserByID(string) *models.User
	GetUserByUsername(string) *models.User
	ExampleClientID() string
}

type userStore struct {
	repository repository.UsersRepository
}

func NewUserStore(issuer string, userRepository repository.UsersRepository) UserStore {
	return userStore{
		repository: userRepository,
	}
}

// ExampleClientID is only used in the example server
func (u userStore) ExampleClientID() string {
	return "service"
}

func (u userStore) GetUserByID(id string) *models.User {
	result, err := u.repository.FindOneById(id)

	if err != nil {
		log.Print(err)
	}

	return result
}

func (u userStore) GetUserByUsername(username string) *models.User {
	result, err := u.repository.FindOneByUserName(username)

	if err != nil {
		log.Print(err)
	}

	return result
}
