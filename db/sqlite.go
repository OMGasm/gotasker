package db

import (
	. "aue.io/tasker/tasks"
	// "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

type SqliteDB struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (self *SqliteDB) Init(dbName string, logger *slog.Logger) error {
	schema := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		created DATETIME,
		date DATETIME
		);`

	var err error
	self.logger = logger
	self.db, err = sqlx.Connect("sqlite3", dbName)
	if err != nil {
		return err
	}
	_, err = self.db.Exec(schema)
	return err
}

func (self SqliteDB) Close() {
	self.db.Close()
}

// a limit of 0 returns all found tasks
func (self SqliteDB) GetTasks(limit int) (tasks []Task, err error) {
	tasks = []Task{}
	query := "SELECT id, title, date FROM tasks"

	if limit <= 0 {
		err = self.db.Select(&tasks, query)
	} else {
		err = self.db.Select(&tasks, query+" LIMIT ?", limit)
	}

	return tasks, err
}

func (self SqliteDB) AddTask(task Task) error {
	query := `INSERT INTO tasks(title, date) VALUES (?, ?)`
	_, err := self.db.Exec(query, task.Title, task.Date)
	return err
}

func (self SqliteDB) DeleteTask(taskId int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := self.db.Exec(query, taskId)
	return err
}
