package main

import (
	_ "database/sql"
	"fmt"
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

func defaultHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprint(response, "Not Found")
		return
	}
	fmt.Fprint(response, "hello world")
}

func main() {
	logger := createLogger()
	slog.SetDefault(logger)

	http.HandleFunc("/", defaultHandler)

	logger.Error("HTTP server has crashed", http.ListenAndServe(":8080", nil))
}
