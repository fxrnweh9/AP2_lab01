package domain

import "context"

type EmailProvider interface {
	SendEmail(
		ctx context.Context,
		to string,
		body string,
	) error
}
