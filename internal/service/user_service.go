package service

import (
	"errors"

	"github.com/fire9900/golang-forum/internal/model"
	"github.com/fire9900/golang-forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user *model.User) error {
	// Проверяем, не занят ли username
	if _, err := s.repo.GetByUsername(user.Username); err == nil {
		return ErrUsernameAlreadyTaken
	}

	// Проверяем, не занят ли email
	if _, err := s.repo.GetByEmail(user.Email); err == nil {
		return ErrEmailAlreadyTaken
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *UserService) Login(login *model.UserLogin) (*model.User, error) {
	user, err := s.repo.GetByEmail(login.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

func (s *UserService) UpdatePassword(id int64, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(id, string(hashedPassword))
}
