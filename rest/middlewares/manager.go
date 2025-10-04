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

	for _, middleware := range middlewares {
		n = middleware(n)
	}

	for _, globalMiddleware := range mngr.globalMiddlewares {
		n = globalMiddleware(n)
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
