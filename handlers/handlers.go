package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	. "aue.io/tasker/db"
	. "aue.io/tasker/tasks"
	"aue.io/tasker/templates"
)

var db DB
var log *slog.Logger

func logHTTP(logger *slog.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			logger.Info("request", "method", req.Method, "path", req.URL.Path)
			handler(res, req)
		}
	}
}

func InitHandlers(logger *slog.Logger, database DB) {
	log = logger
	db = database
}

func RegisterHandlers() {
	log := logHTTP(log)

	http.Handle("GET /", http.FileServer(http.Dir("./public")))
	http.HandleFunc("GET /{$}", log(indexPage))
	http.HandleFunc("GET /tasks", log(taskList))
	http.HandleFunc("DELETE /delete/{id}", log(deleteTask))
	http.HandleFunc("POST /insert", log(insertTask))
	http.HandleFunc("GET /insert", log(insertTask_get))
}

func notFound(reason string, res http.ResponseWriter, req *http.Request) {
	slog.Info(reason, "method", req.Method, "path", req.URL.Path)
	err := templates.E404(req.URL.Path)
	res.WriteHeader(http.StatusNotFound)
	err.Render(req.Context(), res)
}

func indexPage(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.GetTasks(0)
	if err != nil {
		templates.E404("db error?" + err.Error())
	}
	index := templates.Index(tasks)
	index.Render(req.Context(), res)
}

func taskList(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.GetTasks(0)
	if err != nil {
		templates.E404("db error?")
	}
	taskList := templates.TaskList(tasks)
	taskList.Render(req.Context(), res)
}

func insertTask_get(res http.ResponseWriter, req *http.Request) {
	form := templates.Create(false)
	form.Render(req.Context(), res)
}

func insertTask(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Warn("Invalid form data?", "Form error", err)
	}
	var task Task
	task.Title = req.PostFormValue("title")
	task.Date = req.PostFormValue("date")
	if task.Title == "" {
		page := templates.Invalid("did you forget to add a title?")
		page.Render(req.Context(), res)
		return
	}
	db.AddTask(task)
	http.Redirect(res, req, "/", http.StatusFound)
}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 0)
	if err != nil {
		slog.Warn("Invalid task id?", "Task id", req.PathValue("id"))
		res.WriteHeader(http.StatusBadRequest)
		errPage := templates.Invalid("bad id?")
		errPage.Render(req.Context(), res)
		return
	}
	db.DeleteTask(int(id))
	res.WriteHeader(http.StatusOK)
}
