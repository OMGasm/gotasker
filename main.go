package main

import (
	. "aue.io/tasker/handlers"
	. "aue.io/tasker/logging"
	_ "database/sql"
	"log/slog"
	"net/http"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	host, port := "127.0.0.1", "8080"

	logger := CreateLogger()
	slog.SetDefault(logger)
	logHTTP := logger.WithGroup("HTTP")
	RegisterHandlers(logHTTP)

	addr := host + ":" + port
	logHTTP.Info("Starting server on " + addr)
	logHTTP.Error("Server has crashed", http.ListenAndServe(":8080", nil))
}
