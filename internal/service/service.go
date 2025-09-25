package service

import "github.com/Xapsiel/bpla_dashboard/internal/config"

type Repository interface{}

type Service struct {
	*UserService
}

func New(repo Repository, cfg config.OidcConfig) Service {
	return Service{
		UserService: NewUserService(repo, cfg),
	}
}
