package service

import (
	"context"

	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/logging"

	modelCategory "crm/app/internal/domain/category/model"
)

type storage interface {
	Create(ctx context.Context, req modelCategory.CreateCategoryInput) error
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error {
	logging.L(ctx).Debug("CreateCategory called service")

	err := s.storage.Create(ctx, req)
	if err != nil {
		return errors.New("service.CreateCategory")
	}

	return nil
}
