package telemetry

import (
	"context"
	"github.com/DataDog/datadog-go/statsd"
	"sync"
	"time"
)

var Discard = &Telemetry{cli: &nopClient{}}

type Telemetry struct {
	cli     client
	log     logger
	reports []reporter
	wg      sync.WaitGroup
	cancel  func()
}

// New telemetry reporter
func New(log logger, addr, prefix string, tags []string) (*Telemetry, error) {
	if addr == "" {
		return Discard, nil
	}

	cli, err := statsd.New(addr, statsd.WithAsyncUDS(), statsd.Buffered(), statsd.WithNamespace(prefix), statsd.WithTags(tags))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Telemetry{
		cli:     cli,
		log:     log,
		cancel:  cancel,
		reports: []reporter{defaultReporter},
	}

	m.run(ctx)

	return m, nil
}

func (t *Telemetry) run(ctx context.Context) {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, rep := range t.reports {
					rep(ctx, t)
				}
			}
		}
	}()
}

// Stop reporting routine
func (t *Telemetry) Stop() {
	if t.cancel == nil {
		return
	}

	t.cancel()
	t.wg.Wait()
}

// Report adds reporter which would run periodically to report certain metrics
// For example, reporter can measure number of active go routines and report it as gauge.
func (t *Telemetry) Report(rep func(context.Context, *Telemetry)) {
	t.reports = append(t.reports, rep)
}

// Timing send statistics about how log something took
// Based on this metric statsd creates multiple metrics, which include: min, max, avg, medium, 95percentile and count.
func (t *Telemetry) Timing(ctx context.Context, name string, duration time.Duration, tags []string) {
	tags = append(TagsFromContext(ctx), tags...)

	if err := t.cli.Timing(name, duration, tags, 1); err != nil {
		t.log.Warning("An error occurred while reporting metric", map[string]interface{}{
			"error":  err.Error(),
			"metric": name,
		})
	}
}

// Gauge sends absolute value of something, for example number of running go routines.
func (t *Telemetry) Gauge(ctx context.Context, name string, val float64, tags []string) {
	tags = append(TagsFromContext(ctx), tags...)

	if err := t.cli.Gauge(name, val, tags, 1); err != nil {
		t.log.Warning("An error occurred while reporting metric", map[string]interface{}{
			"error":  err.Error(),
			"metric": name,
		})
	}
}

// Incr increases counter by one.
func (t *Telemetry) Incr(ctx context.Context, name string, tags []string) {
	tags = append(TagsFromContext(ctx), tags...)

	if err := t.cli.Incr(name, tags, 1); err != nil {
		t.log.Warning("An error occurred while reporting metric", map[string]interface{}{
			"error":  err.Error(),
			"metric": name,
		})
	}
}

// IncrBy increases counter by given value.
func (t *Telemetry) IncrBy(ctx context.Context, name string, val int64, tags []string) {
	tags = append(TagsFromContext(ctx), tags...)

	if err := t.cli.Count(name, val, tags, 1); err != nil {
		t.log.Warning("An error occurred while reporting metric", map[string]interface{}{
			"error":  err.Error(),
			"metric": name,
		})
	}
}
