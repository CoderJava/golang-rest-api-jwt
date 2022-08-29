package service

import (
	"errors"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/entity"
	"golang-rest-api-jwt/repository"
	"golang-rest-api-jwt/response"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(registerRequest dto.RegisterRequest) (*response.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(registerRequest dto.RegisterRequest) (*response.UserResponse, error) {
	user, err := s.userRepository.FindByEmail(registerRequest.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = entity.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	}

	user, _ = s.userRepository.InsertUser(user)
	result := response.NewUserResponse(user)
	return result, nil
}
