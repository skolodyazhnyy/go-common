//+build go1.11

package telemetry

import (
	"context"
	"database/sql"
)

// DBStats creates reporter for database metrics
func DBStats(db *sql.DB, tags []string) func(ctx context.Context, t *Telemetry) {
	return func(ctx context.Context, t *Telemetry) {
		stats := db.Stats()
		t.Gauge(ctx, "db_open_conn", float64(stats.OpenConnections), tags)
		t.Gauge(ctx, "db_idle_conn", float64(stats.Idle), tags)
		t.Gauge(ctx, "db_used_conn", float64(stats.InUse), tags)
	}
}
