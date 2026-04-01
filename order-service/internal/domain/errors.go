package domain

import "errors"

var (
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrCannotCancelOrder = errors.New("cannot cancel non-pending order")
)
