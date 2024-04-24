package db

import (
	. "aue.io/tasker/tasks"
	"slices"
)

type FakeDB struct {
	tasks []Task // totally ACID compliant ;)
	taskN int
}

func (self FakeDB) Init() error {
	self.AddTask(NewTask("A task to foo", "2024-04-03 15:48"))
	return nil
}

func (self FakeDB) AllTasks() ([]Task, error) {
	return self.tasks, nil
}

func (self FakeDB) FirstNTasks(n int) ([]Task, error) {
	return self.tasks[:n], nil
}

func (self FakeDB) AddTask(task Task) error {
	task.Id = self.taskN
	self.tasks = slices.Insert(self.tasks, len(self.tasks), task)
	self.taskN++
	return nil
}

func (self FakeDB) DeleteTask(id int) error {
	self.tasks = slices.DeleteFunc(self.tasks, func(e Task) bool {
		return e.Id == id
	})
	return nil
}
