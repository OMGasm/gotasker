package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	. "aue.io/tasker/db"
	. "aue.io/tasker/tasks"
	"aue.io/tasker/templates"
)

var db *TaskDB

func defaultHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		notFound("Not found", res, req)
		return
	}
	index := templates.Index(db.AllTasks())
	index.Render(req.Context(), res)
}

func logHTTP(logger *slog.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("request", "method", req.Method, "path", req.URL.Path)
		handler(res, req)
	}
}

func RegisterHandlers(logger *slog.Logger) {
	db = NewDB()

	handlers := map[string]http.HandlerFunc{
		"/":            defaultHandler,
		"POST /delete": deleteTask,
		"POST /insert": insertTask,
		"GET /insert":  insertTask_get,
	}
	for path, handler := range handlers {
		http.HandleFunc(path, logHTTP(logger, handler))
	}
}

func notFound(reason string, res http.ResponseWriter, req *http.Request) {
	slog.Info(reason, "method", req.Method, "path", req.URL.Path)
	err := templates.E404(req.URL.Path)
	res.WriteHeader(http.StatusNotFound)
	err.Render(req.Context(), res)
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
	id, err := strconv.ParseInt(req.PostFormValue("id"), 10, 0)
	if err != nil {
		slog.Warn("Invalid task id?", "Form task id", req.PostFormValue("id"))
		res.WriteHeader(http.StatusBadRequest)
		errPage := templates.Invalid("bad id?")
		errPage.Render(req.Context(), res)
		return
	}
	db.DeleteTask(int(id))
	http.Redirect(res, req, "/", http.StatusFound)
}
