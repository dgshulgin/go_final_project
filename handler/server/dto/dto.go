package dto

// DTO для сообщения об успехе
type Ok struct{}

// DTO для сообщения об ошибке
type Error struct {
	Error string `json:"error"`
}

// DTO для входящего запроса
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// DTO для списка задач
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// DTO для идентификатора задачи
type Id struct {
	Id uint `json:"id"`
}

// пароль пользователя
type Auth struct {
	Password string `json:"password,omitempty"`
}

// JWT
type JWT struct {
	Token string `json:"token"`
}
