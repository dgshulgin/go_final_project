package entity

// Представление задачи в виде записи в БД
type Task struct {
	TaskId  uint   `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}
