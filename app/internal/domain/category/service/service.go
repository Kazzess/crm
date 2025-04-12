package service

import (
	"context"

	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/logging"

	modelCategory "crm/app/internal/domain/category/model"
)

type storage interface {
	CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error
}

type Service struct {
	storage storage
}

func NewStorage(storage storage) Service {
	return Service{
		storage: storage,
	}
}

func (s *Service) CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error {
	logging.L(ctx).Debug("CreateCategory called service")

	err := s.storage.CreateCategory(ctx, req)
	if err != nil {
		return errors.New("service.CreateCategory")
	}

	return nil
}
