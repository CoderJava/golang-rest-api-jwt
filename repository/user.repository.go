package repository

import (
	"fmt"
	"golang-rest-api-jwt/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) (entity.User, error)

	FindByEmail(email string) (entity.User, error)

	FindByUserID(userID string) (entity.User, error)

	UpdateUser(user entity.User) (entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func hashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

func (r *userRepository) InsertUser(user entity.User) (entity.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	r.db.Save(&user)
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	if result := r.db.Where("email = ?", email).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *userRepository) FindByUserID(userID string) (entity.User, error) {
	var user entity.User
	if result := r.db.Where("id = ?", userID).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user entity.User) (entity.User, error) {
	if user.Password == "" {
		var tempUser entity.User
		r.db.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	} else {
		user.Password = hashAndSalt([]byte(user.Password))
	}
	r.db.Save(&user)
	return user, nil
}
