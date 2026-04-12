package services

import (
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]*db.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (*db.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) CreateUser(user *db.User) (*db.User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(id int, user *db.User) (*db.User, error) {
	return s.repo.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
