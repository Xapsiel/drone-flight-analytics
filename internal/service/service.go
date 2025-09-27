package service

import (
	"context"

	"github.com/Xapsiel/bpla_dashboard/internal/config"
	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

type Repository interface {
	SaveMessage(ctx context.Context, mes *model.ParsedMessage) error
}

type Service struct {
	*UserService
	*ParserService
}

func New(repo Repository, cfg config.OidcConfig) Service {
	return Service{
		UserService:   NewUserService(repo, cfg),
		ParserService: NewParserService(repo),
	}
}
