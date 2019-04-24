package telemetry

import (
	"context"
	"io/ioutil"
	"runtime"
)

type reporter func(context.Context, *Telemetry)

func defaultReporter(ctx context.Context, t *Telemetry) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)

	fd, _ := ioutil.ReadDir("/proc/self/fd/")

	t.Gauge(ctx, "mem.allocated", float64(stat.Alloc), nil)
	t.Gauge(ctx, "go_routines", float64(runtime.NumGoroutine()), nil)
	t.Gauge(ctx, "open_fd", float64(len(fd)), nil)
}
