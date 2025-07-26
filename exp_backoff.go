package try

import (
	"errors"
	"time"
)

var ErrAttemptsExceeded = errors.New("attempts exceeded")

var sleep = time.Sleep

func WithExpBackoff(initialDelay, maxDelay time.Duration, maxAttempts int) func(f func() error) error {
	initialDelay = max(initialDelay, 1*time.Millisecond)
	maxDelay = max(maxDelay, initialDelay, 1*time.Millisecond)
	maxAttempts = max(maxAttempts, 1)

	return func(f func() error) error {
		delay := initialDelay
		var err error
		for i := 0; i < maxAttempts; i++ {
			err = f()
			if err == nil {
				return nil
			}

			sleep(delay)
			delay = min(delay*2, maxDelay)
		}

		return errors.Join(ErrAttemptsExceeded, err)
	}
}
