package dto

// DTO сообщения об успехе
type Ok struct{}

// DTO для сообщения об ошибке
type Error struct {
	Error string `json:"error"`
}

// DTO входящего запроса
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// DTO списка задач
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// DTO идентификатора задачи
type Id struct {
	Id uint `json:"id"`
}

// DTO пароль пользователя
type Auth struct {
	Password string `json:"password,omitempty"`
}

// DTO токена
type JWT struct {
	Token string `json:"token"`
}

type RepeatCons struct {
	Date   string `json:"date"`
	Now    string `json:"now"`
	Repeat string `json:"repeat"`
}
