package routing

type Middleware func(EventHandler) EventHandler
