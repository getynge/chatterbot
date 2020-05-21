package routing

// Middleware is a type used to describe any function that can apply logic that runs either before or after the target
// event handler in a given route.
type Middleware func(EventHandler) EventHandler
