package domain

import "errors"

var (
	ErrPaymentLimitExceeded = errors.New("amount exceeds transaction limit")
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
	ErrCannotCancelOrder    = errors.New("cannot cancel non-pending order")
	ErrServiceUnavailable   = errors.New("payment service unavailable")
)
