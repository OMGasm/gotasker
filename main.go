package main

import (
	"aue.io/tasker/templates"
	_ "database/sql"
	"log/slog"
	"net/http"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

func logHTTP(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		slog.Info("request", "method", req.Method, "path", req.URL.Path)
		handler(res, req)
	}
}

func main() {
	host, port := "127.0.0.1", "8080"

	logger := createLogger()
	slog.SetDefault(logger)
	http.HandleFunc("/", logHTTP(defaultHandler))

	addr := host + ":" + port
	logger.Info("Starting HTTP server on " + addr)
	logger.Error("HTTP server has crashed", http.ListenAndServe(":8080", nil))
}
