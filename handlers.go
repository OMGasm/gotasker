package main

import (
	"aue.io/tasker/templates"
	"log/slog"
	"net/http"
)

func defaultHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		slog.Info("Not found", "method", req.Method, "path", req.URL.Path)
		err := templates.E404(req.URL.Path)
		res.WriteHeader(http.StatusNotFound)
		err.Render(req.Context(), res)
		return
	}
	index := templates.Index()
	index.Render(req.Context(), res)
}

func logHTTP(logger *slog.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("request", "method", req.Method, "path", req.URL.Path)
		handler(res, req)
	}
}

func registerHandlers(logger *slog.Logger) {
	handlers := map[string]http.HandlerFunc{
		"/": defaultHandler,
	}
	for path, handler := range handlers {
		http.HandleFunc(path, logHTTP(logger, handler))
	}
}
