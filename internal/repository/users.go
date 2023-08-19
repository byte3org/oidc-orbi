package repository

import (
	"github.com/byte3org/oidc-orbi/internal/infrastructure"
	"github.com/byte3org/oidc-orbi/internal/models"
	"github.com/byte3org/oidc-orbi/internal/utils"
)

type UsersRepository struct {
	infrastructure.Database
}

func NewUsersRepository(db infrastructure.Database) UsersRepository {
	return UsersRepository{db}
}

func (s UsersRepository) Save(user *models.User) error {
	return s.DB.Save(&user).Error
}

func (s UsersRepository) Create(user *models.User) error {

	if user.Password != "" {
		hash, err := utils.MakePassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = hash
	}

	return s.DB.Create(&user).Error
}

func (s UsersRepository) FindOneById(userId string) (*models.User, error) {
	user := models.User{}
	if err := s.DB.Model(&models.User{}).Where("id = ? ", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UsersRepository) FindOneByUserName(userName string) (*models.User, error) {
	user := models.User{}

	if err := s.DB.Model(&models.User{}).Where("username = ? ", userName).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
