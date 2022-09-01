package service

import (
	"errors"
	"fmt"
	"golang-rest-api-jwt/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) error
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) VerifyCredential(email string, password string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	isValidPassword := comparePassword(user.Password, []byte(password))
	if !isValidPassword {
		return errors.New("failed to login. Please check your credential")
	}

	return nil
}

func comparePassword(hasedPassword string, planPasword []byte) bool {
	byteHash := []byte(hasedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, planPasword)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
