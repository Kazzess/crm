package category

import (
	"context"

	gRPCCategoryService "github.com/Kazzess/contracts/gen/go/category_service/v1"
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
