package rpc

// Handler handles RPC request and returns result or an error.
// It is used to process requests received by RPC server.
type Handler interface {
	Handle(Request) (interface{}, error)
}

// HandlerFunc is an RPC handler described as a function.
// It is useful when you don't want to create a structure for handler and just
// want to define handler's logic.
//
//     handler := HandlerFunc(func(Request) (interface{], error) {
//         ...
//     })
//
type HandlerFunc func(Request) (interface{}, error)

// Handle calls handling function.
// This method implements Handler interface.
func (hf HandlerFunc) Handle(req Request) (interface{}, error) {
	return hf(req)
}

// Middleware extends handler with additional capabilities.
//
// An example of logging middleware:
//
//    // Middleware constructor
//    func LoggingMiddleware(logger log.Logger) Middleware {
//
//        // Middleware itself, it takes handler and returns extended handler
//        return func(h Handler) Handler {
//
//            // This function is a handler which logs request and calls "decorated" handler
//            return HandlerFunc(func(req Request) (interface{}, error) {
//                logger.Printf("Calling %v...\n", req.Method)
//                return h(req)
//           })
//        }
//    }
//
// Now we can decorate handler with this middleware:
//
//    h := Decorate(h, LoggingMiddleware(logger))
//
// Middleware can be used to execute some code before or after handler, or even take a decision
// to interrupt request and return immediately, for example authorization middleware.
//
type Middleware func(Handler) Handler

// Decorate extends handler by wrapping it into middleware. Middleware allows
// to execute some logic before or after handler is executed, or even to prevent
// handler from being executed.
//
// For example, security middleware can verify token in the RPC request before
// executing actual handler. Or, metrics middleware could measure latency of the
// handler execution.
//
// Decorate will take handler and middleware and return a new, decorated, handler.
//
//    handler = Decorate(
//        rpc.HandlerFunc(func(rpc.Request) {
//            ...
//        }),
//        NewMetricsMiddleware(),
//        NewSecurityMiddleware(),
//        NewLoggingMiddleware(),
//    )
//
// Middleware is executed in reversed order, last in the list will be executed first.
func Decorate(h Handler, mw ...Middleware) Handler {
	for _, m := range mw {
		h = m(h)
	}

	return h
}
