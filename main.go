package main

import (
	_ "database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func createLogger() (*slog.Logger, *log.Logger) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	sink := io.MultiWriter(file, stdout)
	handler := slog.NewJSONHandler(sink, nil)
	logger := slog.New(handler)
	basicLogger := slog.NewLogLogger(handler, slog.LevelInfo)
	return logger, basicLogger
}

func defaultHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprint(response, "Not Found")
		return
	}
	fmt.Fprint(response, "hello world")
}

func main() {
	logger, basicLogger := createLogger()
	slog.SetDefault(logger)

	http.HandleFunc("/", defaultHandler)

	basicLogger.Fatal(http.ListenAndServe(":8080", nil))
}
