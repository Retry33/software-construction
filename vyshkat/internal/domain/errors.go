package domain

import "errors"

var (
	ErrEmptyName             = errors.New("name must not be empty")
	ErrNegativeEnergy        = errors.New("daily energy consumption must be a non-negative number")
	ErrInvalidSimplicity     = errors.New("simplicity must be between 1 and 10")
	ErrEmptyNumber           = errors.New("inventory number must not be empty")
	ErrNumberAlreadyAssigned = errors.New("inventory number is already assigned")
)
