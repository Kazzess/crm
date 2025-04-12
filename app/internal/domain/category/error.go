package category

import (
	"github.com/Kazzess/libraries/errors"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("already exists")
)

// -------------------------------------- Errors and constants from storage  --------------------------------------

const (
	CategoryIDPkConstraint = "category_id_pk"
)

var (
	ErrViolatesConstraintCategoryIDPk = errors.New("violates constraint category_id_pk")
)
