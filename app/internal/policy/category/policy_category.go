package category

import (
	"context"
	"errors"
	"time"

	modelCategory "crm/app/internal/domain/category/model"
)

func (p *Policy) CreateCategory(ctx context.Context, name string) error {
	// TODO add base policy: UUID generator, clock generator.
	time := time.Now()
	request := modelCategory.ConstructorCreateCategoryInput(
		name,
		false,
		time,
		time,
	)

	err := p.service.CreateCategory(ctx, request)
	if err != nil {
		return errors.New("service.CreateCategory")
	}

	return nil
}
