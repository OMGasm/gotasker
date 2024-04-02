package tasks

type Task struct {
	title string
	date  string
}

func NewTask(title, date string) Task {
	return Task{
		title,
		date,
	}
}

func (self Task) Title() string {
	return self.title
}

func (self Task) Date() string {
	return self.date
}
