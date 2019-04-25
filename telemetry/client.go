package telemetry

import "time"

type client interface {
	Timing(name string, value time.Duration, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Incr(name string, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

type nopClient struct {
}

func (nopClient) Timing(name string, value time.Duration, tags []string, rate float64) error {
	return nil
}

func (nopClient) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (nopClient) Incr(name string, tags []string, rate float64) error {
	return nil
}

func (nopClient) Count(name string, value int64, tags []string, rate float64) error {
	return nil
}
