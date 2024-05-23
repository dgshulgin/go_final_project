package rules

import (
	"time"
)

// Реализация интерфейса RepeatRule для правила повторения y.
// Вычисляет ближайшую дату события, если задача выполняется ежегодно.
// Выполнение задачи переносится на один год вперед.
// Правила формирования условия repeat:
// Формат: "y"

type RepeatYears struct {
}

func NewRepeatYears() *RepeatYears { return &RepeatYears{} }

func (rd RepeatYears) Reset() {}

func (rd RepeatYears) Validate(repeat string) error {
	return nil
}
func (rd RepeatYears) Apply(start time.Time) time.Time {
	return start.AddDate(1, 0, 0)
}
