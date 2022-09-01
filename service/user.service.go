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

	FindUserByEmail(email string) (*response.UserResponse, error)

	FindUserByID(userID string) (*response.UserResponse, error)

	UpdateUser(updateUserRequest dto.UpdateUserRequest) (*response.UserResponse, error)
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

func (s *userService) FindUserByEmail(email string) (*response.UserResponse, error) {
	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	userResponse := response.NewUserResponse(user)
	return userResponse, nil
}

func (s *userService) FindUserByID(userID string) (*response.UserResponse, error) {
	user, err := s.userRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	userResponse := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return &userResponse, nil
}

func (s *userService) UpdateUser(updateUserRequest dto.UpdateUserRequest) (*response.UserResponse, error) {
	user := entity.User{
		ID:    updateUserRequest.ID,
		Name:  updateUserRequest.Name,
		Email: updateUserRequest.Email,
	}

	user, err := s.userRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	response := response.NewUserResponse(user)
	return response, nil
}
