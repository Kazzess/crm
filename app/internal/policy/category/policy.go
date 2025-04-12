package category

import (
	"context"

	modelCategory "crm/app/internal/domain/category/model"

	"crm/app/internal/policy"
)

type Service interface {
	CreateCategory(ctx context.Context, req modelCategory.CreateCategoryInput) error
}

type Policy struct {
	*policy.BasePolicy
	service Service
}

func NewPolicy(
	basePolicy *policy.BasePolicy,
	service Service,
) *Policy {
	return &Policy{
		BasePolicy: basePolicy,
		service:    service,
	}
}
