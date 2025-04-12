package category

import (
	"context"

	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/logging"

	gRPCCategoryService "github.com/Kazzess/contracts/gen/go/category_service/v1"
)

//  ------------------------------------------------ Action with Category -----------------------------------------

// CreateCategory category for CRM.
func (c *Controller) CreateCategory(
	ctx context.Context,
	req *gRPCCategoryService.CreateCategoryRequest,
) (*gRPCCategoryService.CreateCategoryResponse, error) {
	logging.L(ctx).Debug("CreateCategory")

	err := c.policy.CreateCategory(ctx, req.GetName())
	if err != nil {
		return nil, errors.Wrap(err, "policyCreateCategory")
	}

	response := gRPCCategoryService.CreateCategoryResponse{}

	return &response, nil
}
