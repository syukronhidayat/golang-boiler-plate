package repositories

import (
	"context"
	"fmt"
	"golang-boiler-plate/config"
	"golang-boiler-plate/models"
	"golang-boiler-plate/utils/types"

	"github.com/go-resty/resty/v2"
)

type UsersRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type userRepository struct {
	client *resty.Client
}

func NewUserRepository() UsersRepository {
	return &userRepository{
		client: resty.New(),
	}
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)

	headers := types.ObjStr{
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
		"Authorization":        fmt.Sprintf("Bearer %s", config.GithubAccessToken),
	}
	_, err := u.client.R().SetHeaders(headers).SetResult(user).Get(fmt.Sprintf("%s/users/%s", config.GithubApiBaseUrl, username))
	if err != nil {
		return nil, err
	}

	return user, nil
}
