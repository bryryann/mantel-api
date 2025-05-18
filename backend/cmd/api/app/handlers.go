package app

import "net/http"

func (a *App) RegisterHandler(method, path string, handler http.HandlerFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.routes = append(a.routes, Route{
		Path:    path,
		Method:  method,
		Handler: handler,
	})
}

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodGet, path, handler)
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodPost, path, handler)
}

func (a *App) Put(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodPut, path, handler)
}

func (a *App) Delete(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodDelete, path, handler)
}
