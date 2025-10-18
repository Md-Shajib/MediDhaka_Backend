package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (manager *Manager) Use(middlewares ...Middleware) {
	manager.globalMiddlewares = append(manager.globalMiddlewares, middlewares...)
}

func (mngr *Manager) With(next http.Handler, middlewares ...Middleware) http.Handler {
	n := next

	// Apply local middlewares first (inner)
	for i := len(middlewares) - 1; i >= 0; i-- {
		n = middlewares[i](n)
	}

	// Then apply global middlewares (outer)
	for i := len(mngr.globalMiddlewares) - 1; i >= 0; i-- {
		n = mngr.globalMiddlewares[i](n)
	}

	return n
}

func (mngr *Manager) WrapMux(handler http.Handler, middlewares ...Middleware) http.Handler {
	h := handler
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
