package db

import (
	. "aue.io/tasker/tasks"
	"log/slog"
)

type DB interface {
	Init(dbName string, logger *slog.Logger) error
	Close()
	GetTasks(n int) ([]Task, error)
	AddTask(Task) error
	DeleteTask(taskId int) error
}
