package rules

import (
	"errors"
	"time"
)

var (
	ErrRuleEmptyKey   = errors.New("отсутствует значение для ключа")
	ErrRuleInvalidKey = errors.New("некорректное значение ключа")
)

// Правила повторения задачи. Каждый ключ d,y,w,m и дополнительный ключ e (для пустых правил) обрабатывается собственной реализацией интерфейса RepeatRule.
type RepeatRule interface {
	// Выполняет проверку корректности условия repeat
	Validate(repeat string) error
	// Применяет условие repeat к начальной дата, возвращает дату следующего
	// наступления события
	Apply(start time.Time) time.Time
	// Служебный метод для очистки внутренного состояния экземпляра RepeatRule.
	// Обязателен для применения, вызывается через defer
	Reset()
}
