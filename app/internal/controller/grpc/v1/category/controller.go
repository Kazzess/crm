package category

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Kazzess/Contracts/gen/go/category_service/v1;pb_category_service_v1"
)

//  ------------------------------------------------ Action with Category -----------------------------------------

// CreateCategory category for CRM.
func (c *Controller) CreateCategory(
	ctx context.Context,
	req *gRPCCategoryService.CreateCategoryRequest,
) (*gRPCCategoryService.CreateCategoryResponse, error) {
	err := c.policy.CreateCategory(ctx, req.GetName())

	if err != nil {
		return nil, errors.Wrap(err, "policyCreateCategory")
	}

	response := gRPCCategoryService.CreateCategoryResponse

	return &response, nil
}
