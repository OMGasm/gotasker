package main

import (
	_ "database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func createLogger() *slog.Logger {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	sink := io.MultiWriter(stdout, file)
	handler := slog.NewJSONHandler(sink, nil)
	logger := slog.New(handler)
	return logger
}

func defaultHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		slog.Info("Not found", "method", req.Method, "path", req.URL.Path)
		err := e404(req.URL.Path)
		res.WriteHeader(http.StatusNotFound)
		err.Render(req.Context(), res)
		return
	}
	index := index()
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
