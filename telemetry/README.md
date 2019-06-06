# Telemetry

Telemetry reporter package to send metrics to datadog.

## Usage

Create telemetry reporter using `New` constructor. It takes, logger instance, Statsd address, application prefix (must end with a dot) and list of global tags attached to all metrics.

In case statsd address is empty, constructor will return `Discard` telemetry reporter.

```go
meter, err := telemetry.New(log.Channel("telemetry"), config.StatsdAddr, "my-app.", nil)
if err != nil {
	log.Fatal("Failed to connect to statsd:", err)
}

defer meter.Stop()
```

When new telemetry structure is created, it automatically starts background routine which reports 
number of running go routines, memory usage and number of open file descriptors. 

Use `Timing`, `Gauge`, `Incr` and `IncrBy` to report metrics in your application.

```go
meter.Timing(ctx, "rpc_request", time.Since(start), []string{
	"method:"+req.Method,
	"client:"+req.Client,
})

meter.IncrBy(ctx, "processed_messages", processed, nil)
```

## Overview

This section provides overview of telemetry API.

### Timing

Timing reports how long something took, it automatically expanded into 6 metrics by datadog: min, max, avg, median, 95percentile and count.

```go
func SelectCount(ctx context.Context, db *sql.DB) (int, error) {
    start := time.Now()
    defer func() {
    	meter.Timing(ctx, "db_select", time.Since(start), []string{"select:user_count"})
    }()
    
    return db.Select(...)
}
```

Example above will create 6 metrics in DataDog:

* `db_select.min` minimal amount of time query took
* `db_select.max` maximal amount of time query took
* `db_select.avg` average amount of time query took
* `db_select.median` 50th percentile 
* `db_select.95percentile` 95th percentile 
* `db_select.count` number of times `SelectCount` was invoked (metric was reported) 

### Incr, IncrBy

Incr and IncrBy report number of times something has happened.

```go
meter.IncrBy(ctx, "messages_processed", processed, nil)
```

Incr increases value by one, each time this function is called a package is sent to datadog daemon. It's highly recommended
to batch up counter increases when possible and use IncrBy, this way number of packages sent to datadog daemon can be reduced. 


### Gauge

Gauge reports absolute value of something, for example number of routines running, amount of memory used by application etc.

```go
meter.Gauge(ctx, "go_rountines", runtime.GoRoutines(), nil)
```

Keep in mind, number of go routines, memory usage and number of open file descriptors are already reported by telemetry service. 


### Report
 
Report adds a callback which is triggered periodically to report metrics. This is useful to report Gauge metrics. 

```go
meter.Report(func(ctx context.Context, t *telemetry.Telemetry) {
    meter.Gauge(ctx, "go_rountines", runtime.GoRoutines(), nil)
})
```

Keep in mind, number of go routines, memory usage and number of open file descriptors are already reported by telemetry service.

### Discard

In case you need to inject telemetry in tests, use `telemetry.Discard`. It's a telemetry instance which throws away all metrics.

### AppendContext

You can assign tags to context using `telemetry.AppendContext`. All metrics reported with appended context will include these tags.

```go
ctx = telemetry.AppendContext(ctx, []string{"merchant:luma"})

meter.Incr(ctx, "request_count", []string{"topic:create_order"}) 
// will report metric with tags: "merchant:luma" and "topic:create_order"
```

### DBStats

Database stats (number of open connections and other metrics) can be reported using `telemetry.DBStats` reporter. It takes database connection and list of tags, for example with database name.

```go
meter.Report(telemetry.DBStats(db, []string{"db:default"}))
```

In case application uses multiple connections, we can assign different tags:

```go
meter.Report(telemetry.DBStats(dbWriter, []string{"db:writer"}))
meter.Report(telemetry.DBStats(dbReader, []string{"db:reader"}))
```
