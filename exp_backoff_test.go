package try

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWithExpBackoff(t *testing.T) {
	var errAny = fmt.Errorf("any")

	tests := []struct {
		name         string
		callResults  []error
		initialDelay time.Duration
		maxDelay     time.Duration
		maxAttempts  int
		expectSleeps []time.Duration
		expectErr    error
	}{
		{
			name:         "success on first try",
			callResults:  []error{nil},
			initialDelay: 10 * time.Millisecond,
			maxDelay:     50 * time.Millisecond,
			maxAttempts:  3,
			expectSleeps: []time.Duration{},
			expectErr:    nil,
		},
		{
			name:         "success after 2 retries",
			callResults:  []error{errAny, errAny, nil},
			initialDelay: 10 * time.Millisecond,
			maxDelay:     50 * time.Millisecond,
			maxAttempts:  5,
			expectSleeps: []time.Duration{10 * time.Millisecond, 20 * time.Millisecond},
			expectErr:    nil,
		},
		{
			name:         "fails when max attempts exceeded",
			callResults:  []error{errAny, errAny, errAny, errAny},
			initialDelay: 10 * time.Millisecond,
			maxDelay:     25 * time.Millisecond,
			maxAttempts:  4,
			expectSleeps: []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 25 * time.Millisecond, 25 * time.Millisecond},
			expectErr:    ErrAttemptsExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slept := make([]time.Duration, 0)

			sleep = func(d time.Duration) {
				slept = append(slept, d)
			}

			results := tt.callResults
			retry := WithExpBackoff(tt.initialDelay, tt.maxDelay, tt.maxAttempts)
			err := retry(func() error {
				result := results[0]
				results = results[1:]
				return result
			})

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			}

			assert.EqualValues(t, tt.expectSleeps, slept)
		})
	}
}
