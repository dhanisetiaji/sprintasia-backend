package services

import (
	"gotham/models"
	"gotham/repositories"
)

type IAuthService interface {
	GetUserByEmail(email string) (user models.User, err error)
	Check(email string, password string) (bool, error)
	Register(register models.RegisterUser) (err error)
}

type AuthService struct {
	UserRepository repositories.IUserRepository
}

func (service *AuthService) Check(email string, password string) (bool, error) {
	user, err := service.UserRepository.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user.VerifyPassword(password), err
}

func (service *AuthService) GetUserByEmail(email string) (user models.User, err error) {
	return service.UserRepository.GetUserByEmail(email)
}

func (service *AuthService) Register(register models.RegisterUser) (err error) {
	err = service.UserRepository.Create(&models.User{}, register)
	if err != nil {
		return err
	}
	return nil
}
