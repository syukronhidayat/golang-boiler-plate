package services

import (
	"context"
	"golang-boiler-plate/models"
	"golang-boiler-plate/repositories"
	"golang-boiler-plate/utils/logger"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type userService struct {
	repository repositories.UsersRepository
}

func NewUserService() UserService {
	return &userService{
		repository: repositories.NewUserRepository(),
	}
}

func (u *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := u.repository.GetByUsername(ctx, username)
	if err != nil {
		logger.Ctx(ctx).StackTrace(err).Error("Error fetching user")
		return nil, err
	}

	return user, nil
}
