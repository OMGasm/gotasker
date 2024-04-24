package main

import (
	"aue.io/tasker/db"
	. "aue.io/tasker/handlers"
	. "aue.io/tasker/logging"
	"log/slog"
	"net/http"
)

func main() {
	host, port := "127.0.0.1", "8080"
	dbName := "tasks.db"

	logger := CreateLogger()
	slog.SetDefault(logger)
	logHTTP := logger.WithGroup("HTTP")
	logDB := logger.WithGroup("DB")

	db := new(db.SqliteDB)
	err := db.Init(dbName, logDB)
	if err != nil {
		logDB.Error("Could not connect", "db", dbName, "error", err)
	}

	InitHandlers(logHTTP, db)
	RegisterHandlers()

	addr := host + ":" + port
	logHTTP.Info("Starting server on " + addr)
	logHTTP.Error("Server has crashed", http.ListenAndServe(":8080", nil))
}
