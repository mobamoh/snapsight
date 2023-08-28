package models

import "github.com/mobamoh/snapsight/errors"

var (
	// A common pattern is to add the package as a prefix to the error for context.
	ErrNotFound   = errors.New("models: resource could not be found")
	ErrEmailTaken = errors.New("models: email address is already in use")
)
