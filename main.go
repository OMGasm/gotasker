package main

import (
	_ "database/sql"
	"log/slog"
	"net/http"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	host, port := "127.0.0.1", "8080"

	logger := createLogger()
	slog.SetDefault(logger)
	logHTTP := logger.WithGroup("HTTP")
	registerHandlers(logHTTP)

	addr := host + ":" + port
	logHTTP.Info("Starting server on " + addr)
	logHTTP.Error("Server has crashed", http.ListenAndServe(":8080", nil))
}
