package category

import (
	"context"
)

type policy interface {
	CreateCategory(ctx context.Context, name string) error
}

type Controller struct {
	gRPCCategoryService.UnimplementedCategoryServiceServer
	policy policy
}

func NewController(policy policy) Controller {
	return Controller{
		policy: policy,
	}
}
