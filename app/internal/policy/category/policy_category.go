package category

import (
	"context"

	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/logging"

	modelCategory "crm/app/internal/domain/category/model"
)

func (p *Policy) CreateCategory(ctx context.Context, name string) error {
	logging.L(ctx).Debug("CreateCategory policy")

	request := modelCategory.ConstructorCreateCategoryInput(
		name,
		false,
		p.Now(),
		p.Now(),
	)

	err := p.service.CreateCategory(ctx, request)
	if err != nil {
		return errors.New("service.CreateCategory")
	}

	return nil
}
