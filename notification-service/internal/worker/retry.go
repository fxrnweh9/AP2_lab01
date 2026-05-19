package worker

import (
	"log"
	"time"
)

func WithRetry(maxRetries int, fn func() error) error {
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt == maxRetries-1 {
			break
		}

		backoff := time.Duration(1<<uint(attempt)) * time.Second
		log.Printf("[Retry] Attempt %d failed: %v. Retrying in %s...", attempt+1, err, backoff)
		time.Sleep(backoff)
	}

	return err
}
