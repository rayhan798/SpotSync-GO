package user

import (
	"errors"

	"gorm.io/gorm"
)

// ErrorAlreadyExist
var ErrorAlreadyExist = errors.New("user with this email already exist")

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	CheckAdminExists() (bool, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) CheckAdminExists() (bool, error) {
	var count int64

	err := r.db.Model(&User{}).Where("role = ?", "admin").Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// CreateUser
func (r *repository) CreateUser(user *User) error {
	var count int64
	r.db.Model(&User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return ErrorAlreadyExist
	}

	result := r.db.Create(user)
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExist
		}
		return result.Error
	}

	return nil
}

// GetUserByEmail
func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User

	result := r.db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}
