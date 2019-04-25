# Logger

Simple structured data logger.

## Usage

Best to create logger first thing in your application. This way any error, even these on very first steps of bootstrapping 
application, can be logged using proper output format.   
 
Logger can print messages in human friendly format, which is useful for debugging or in JSON which is useful for environments
where logs are collected by log collector like Splunk Agent or Logstash. 

Use flags to define format of the message instead of configuration file, because in case configuration can not be loaded your
logger won't know which format to use to print "configuration can not be loaded".

```go
func main() {
	logfmt := flag.String("logfmt", "json", "Logger format (options are: json or text)")
	flag.Parse()
	
	log := logger.New(*logfmt, os.Stdout)
	// logger is ready to be used! 
}
```

Use `Fatal` function to report critical bootstrapping errors in your applications. This call will stop application execution
with non-zero exit code. 

```go
if err := sql.Connect("mysql", *dsn); err != nil {
	log.Fatal("Unable to connect to database:", err)
	// no need to os.Exit here because Fatal will do it for you
}
```

Never use `Fatal` in places other than application bootstrapping. 

Next, inject logger into other services which require logging. Use `Channel` function to differentiate where messages are coming from.

```go
store := user.NewStore(db, log.Channel("user-store"))
```

Don't use `Channel` in places other than application bootstrapping. 

Use `Debug`, `Info`, `Warning` and `Error` methods to log records with different severity level. All these methods accept message and 
map of context parameters. Message should be constant string and all variables should be passed as parameters.

```go
log.Error("User select query has failed", map[string]interface{}{
	"error": err.Error(),
	"user_id": id,
})
``` 

Constant error messages simplify error lookup, filtering and aggregation.

Some parameters have special meaning and printed or processed differently. 

* `error` used to report an error message
* `context` used to extract request related information from context, see "AppendContext" section below
* `channel` used to report a channel, normally should not be used directly but instead through `Channel` constructor
* `level` reserved for internal use
* `message` reserved for internal use
* `time` reserved for internal use
 
You can also inject logger into services which expect standard `log.Logger` interface. Use `NewStdLogger` constructor to turn logger into 
its version compatible with `log.Logger` interface.
 
## Overview

This section provides overview of logger API.

### Debug, Info, Warning, Error

These methods are used to report errors with different level of severity. They all accept message string and list of parameters.

```go
log.Debug("Message is received", map[string]interface{}{
	"message_id": msg.ID,
	"message_topic": msg.Topic,
})
```

### Channel

Channel allows to create new logger from existing one with `channel` parameter preset.

```go
clog = log.Channel("user-store")
clog.Debug("Test", nil) 
// will print: {..., "msg": "Test", "level": "DEBUG", "channel": "user-store"}
```

### With

With allows to create new logger from existing one with predefined parameters.

```go
log = log.With(map[string]interface{}{
	"service": "my-app",
	"version": "12abe2",
})
log.Debug("Test", nil)
// will print: {..., "msg": "Test", "level": "DEBUG", "service": "my-app", "version": "12abe2"}
```

### Fatal

Fatal allows to signal a fatal error in application (typically during bootstrapping) and stop application execution.
                                                                           
```go                                                                      
if err := sql.Connect("mysql", *dsn); err != nil {                         
	log.Fatal("Unable to connect to database:", err)                       
	// no need to os.Exit here because, Fatal will do it for you           
}                                                                          
```                                                                        

### Discard

In case you need to inject logger in tests but don't wish to see logger output, use `logger.Discard`. 
It's a logger instance which throws away all records.

### AppendContext

In some cases you might want to assign some parameters to request and log them with every message relevant for the request.
In this case, you can use `logger.AppendContext` function which would create a copy of `context.Context` with additional logging parameters.
Then, use `context` parameter when logging messages to extract additional parameters from context.

```go
// assign parameters to ctx
ctx = logger.AppendContext(ctx, map[string]interface{}{
	"request_id": id,
})

// tell logger to lookup parameters in context
log.Debug("Test", map[string]interface{}{
	"context": ctx,
})

// will print: {..., "msg": "Test", "level": "DEBUG", "request_id": "..."}
```
