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
		"/":       defaultHandler,
		"/delete": deleteTask,
		"/insert": insertTask,
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

func insertTask(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		form := templates.Create(false)
		form.Render(req.Context(), res)

	case http.MethodPost:
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

	default:
		notFound("Invalid method", res, req)
	}

}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete || req.Method == http.MethodPost {
		id, err := strconv.ParseInt(req.URL.Query().Get("id"), 10, 0)
		if err != nil {
			slog.Warn("Invalid task id?", "Query id", req.URL.Query().Get("id"))
			res.WriteHeader(http.StatusBadRequest)
			errPage := templates.Invalid("bad id?")
			errPage.Render(req.Context(), res)
			return
		}
		db.DeleteTask(int(id))
		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		notFound("Invalid method", res, req)
	}
}
