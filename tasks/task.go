package tasks

type Task struct {
	Id    int `json:"-"`
	Title string
	Date  string
}

func NewTask(title, date string) Task {
	return Task{
		0,
		title,
		date,
	}
}
