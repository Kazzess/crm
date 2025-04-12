package dal

import (
	"github.com/Kazzess/libraries/errors"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
