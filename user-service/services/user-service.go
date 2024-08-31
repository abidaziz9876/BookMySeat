package services

import (
	"errors"

	"github.com/abidaziz9876/user-service/models"
	"github.com/abidaziz9876/user-service/repository"
)

func CreateUser(postgres repository.PostgesDatabase, User models.UserModel) error {
	existingUser, err := postgres.FindUserByEmailOrPhone(User.Email, User.Phone)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}
	result := postgres.Create(&User)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
