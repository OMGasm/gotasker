package logging

import (
	"io"
	"log/slog"
	"os"
)

func CreateLogger() *slog.Logger {
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
