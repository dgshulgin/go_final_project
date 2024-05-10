package rules

import (
	"time"
)

const (
	stepYears = 1
)

// Вычисляет ближайшую дату события, если задача выполняется ежегодно.
// Выполнение задачи переносится на один год вперед.
// Правила формирования условия repeat:
// Формат: "y"

type RepeatYears struct {
}

func NewRepeatYears() *RepeatYears { return &RepeatYears{} }

func (rd RepeatYears) Validate(params []string) error {
	return nil
}
func (rd RepeatYears) Apply(start time.Time) time.Time {
	return start.AddDate(stepYears, 0, 0)
}
