package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"
)

type MockEmailProvider struct{}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

func (m *MockEmailProvider) SendEmail(
	ctx context.Context,
	to string,
	body string,
) error {

	time.Sleep(500 * time.Millisecond)

	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("temporary network error")
	}

	log.Printf(
		"[MockEmail] Successfully sent to %s",
		to,
	)

	return nil
}
