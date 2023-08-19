package repository

import "github.com/byte3org/oidc-orbi/internal/infrastructure"

type UsersRepository struct {
	infrastructure.Database
}

func NewUsersRepository(db infrastructure.Database) UsersRepository {
	return UsersRepository{db}
}

func (s UsersRepository) FindAll(limit int, offset int) (response *UserFindAllResult, err error) {
	var users []models.User
	var count int64

	if err := s.DB.Debug().Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	if err := s.DB.Find(&models.User{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &UserFindAllResult{Users: users, Count: count}, nil
}

func (s UsersRepository) Save(user *models.User) error {
	return s.DB.Save(&user).Error
}

func (s UsersRepository) Create(user *models.User) error {
	return s.DB.Create(&user).Error
}

func (s UsersRepository) FindOneById(userId string) (*models.User, error) {
	user := models.User{}
	if err := s.DB.Model(&models.User{}).Where("id = ? ", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UsersRepository) FindOneByEmail(email string) (*models.User, error) {
	user := models.User{}

	if err := s.DB.Model(&models.User{}).Where("email = ? ", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UsersRepository) FindOneByProviderId(sub string) (*models.User, error) {
	user := models.User{}
	if err := s.DB.Model(&models.User{}).Where("sub = ? ", sub).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UsersRepository) SetUserProvider(userId string, sub string) error {
	return s.DB.Model(&models.User{}).Where("id = ? ", userId).Update("sub", sub).Error
}

func (s UsersRepository) FindAllSpeakers(Speakers []string) (*[]models.Speaker, error) {
	var speakers []models.Speaker

	if err := s.DB.Debug().Where("id IN (?)", Speakers).Find(&speakers).Error; err != nil {
		return nil, err
	}

	return &speakers, nil
}
