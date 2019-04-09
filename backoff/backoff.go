package backoff

import "time"

const defaultMaxBackOffTimeout = time.Minute
const defaultResetAfter = 10 * time.Second

// Backoff handles exponential increases
type Backoff struct {
	last       time.Time
	delay      time.Duration
	maximum    time.Duration
	resetAfter time.Duration
}

// New instances a backoff with default parameters
func New() *Backoff {
	return &Backoff{
		maximum:    defaultMaxBackOffTimeout,
		resetAfter: defaultResetAfter,
	}
}

// NewWithParameters instances a backoff specifying parameters
// Maximum parameter is to specify maximum duration of the delay
// ResetAfter will help you to determine after when the delay is gonna be reset
func NewWithParameters(maximum time.Duration, resetAfter time.Duration) *Backoff {
	return &Backoff{
		maximum:    maximum,
		resetAfter: resetAfter,
	}
}

// Timeout returns backoff delay
func (b *Backoff) Timeout() time.Duration {
	// reset delay if last run took longer than 10 seconds (without backoff delay)
	if time.Since(b.last)-b.delay > b.resetAfter {
		b.delay = 0
	}

	b.last = time.Now()
	b.delay = b.increase(b.delay)

	return b.delay
}

// increase backoff delay
func (b *Backoff) increase(d time.Duration) time.Duration {
	if d <= 0 {
		return time.Second
	}

	if d*2 > b.maximum {
		return b.maximum
	}

	return d * 2
}
