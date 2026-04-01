package domain

import "errors"

var (
	ErrPaymentLimitExceeded = errors.New("amount exceeds transaction limit")
)
