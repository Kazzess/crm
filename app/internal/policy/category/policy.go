package category

import (
	"context"

	modelCategory "crm/app/internal/domain/category/model"
)

type Service interface {
	CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error
	UpdateCategory()
	SearchCategory()
}

type Policy struct {
	service Service
}

func NewPolicy(
	service Service,
) Policy {
	return Policy{
		service: service,
	}
}
