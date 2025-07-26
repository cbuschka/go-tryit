package examples

import (
	"github.io/cbuschka/go-tryit"
	"os"
	"time"
)

func basicExample() error {
	runWithRetry := tryit.WithExpBackoff(1*time.Millisecond,
		300*time.Millisecond, 5)
	err := runWithRetry(func() error {
		return os.Chmod("/nonexistent", 0755)
	})
	return err
}
