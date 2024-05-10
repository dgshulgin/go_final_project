package rules

import (
	"fmt"
	"strconv"
	"time"
)

const (
	minDays = 1
	maxDays = 400
)

// Вычисляет ближайшую дату события, если условие повторения задано как кол-во дней.
// Правила формирования условия repeat:
// Формат: "d <число>", где число находится в пределах 1-400
type RepeatDays struct {
	days int
}

func NewRepeatDays() *RepeatDays { return &RepeatDays{} }

func (rd *RepeatDays) Validate(params []string) error {

	// проверка формата
	if len(params) == 0 {
		return fmt.Errorf("отсутствует значение для ключа d")
	}

	// проверка предельных значений ключа
	days, err := strconv.Atoi(params[0])
	if err != nil {
		return fmt.Errorf("некорректное значение аргумента для ключа d: %s", params[0])
	}
	if days < minDays || days > maxDays {
		return fmt.Errorf("недопустимое количество дней для ключа d %d (%d-%d))", days, minDays, maxDays)
	}

	rd.days = days

	return nil
}
func (rd RepeatDays) Apply(start time.Time) time.Time {
	return start.AddDate(0, 0, rd.days)
}
