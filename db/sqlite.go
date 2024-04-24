package db

import (
	. "aue.io/tasker/tasks"
	"database/sql"
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
	var rows *sql.Rows
	query := "SELECT id, title, date FROM tasks"
	if limit <= 0 {
		rows, err = self.db.Query(query)
	} else {
		rows, err = self.db.Query(query+" LIMIT ?", limit)
	}
	tasks = []Task{}
	for rows.Next() {
		var task Task
		rows.Scan(&task.Id, &task.Title, &task.Date)
		tasks = append(tasks, task)
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
