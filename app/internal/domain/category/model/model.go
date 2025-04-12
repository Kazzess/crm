package model

import "time"

type CreateCategoryInput struct {
	Name      string
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConstructorCreateCategoryInput(
	name string,
	isEnabled bool,
	createdAt, updatedAt time.Time,
) CreateCategoryInput {
	return CreateCategoryInput{
		Name:      name,
		IsEnabled: isEnabled,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
