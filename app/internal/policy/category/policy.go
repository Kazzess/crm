package category

import (
	"context"

	modelCategory "crm/app/internal/domain/category/model"

	"crm/app/internal/policy"
)

type Service interface {
	CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error
	UpdateCategory()
	SearchCategory()
}

type Policy struct {
	policy.BasePolicy
	service Service
}

func NewPolicy(
	service Service,
) Policy {
	return Policy{
		BasePolicy: policy.BasePolicy{},
		service:    service,
	}
}
