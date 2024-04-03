package db

import (
	. "aue.io/tasker/tasks"
	"slices"
)

type TaskDB struct {
	tasks []Task // totally ACID compliant ;)
	taskN int
}

func NewDB() *TaskDB {
	db := TaskDB{}
	db.AddTask(NewTask("A task to foo", "2024-04-03 15:48"))

	return &db
}

func (self *TaskDB) AllTasks() []Task {
	return self.tasks
}

func (self *TaskDB) FirstNTasks(n int) []Task {
	return self.tasks[:n]
}

func (self *TaskDB) AddTask(task Task) {
	task.Id = self.taskN
	self.tasks = slices.Insert(self.tasks, len(self.tasks), task)
	self.taskN++
}

func (self *TaskDB) DeleteTask(id int) {
	self.tasks = slices.DeleteFunc(self.tasks, func(e Task) bool {
		return e.Id == id
	})
}
