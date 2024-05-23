package rules

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	minDays = 1
	maxDays = 400
)

// Реализация интерфейса RepeatRule для правила повторения d.
// Вычисляет ближайшую дату события, если условие повторения задано
// как кол-во дней.
// Правила формирования условия repeat:
// Формат: "d <число>", где число находится в пределах 1-400
type RepeatDays struct {
	days int
}

func NewRepeatDays() *RepeatDays { return &RepeatDays{} }

func (rd *RepeatDays) Reset() { rd.days = 0 }

func (rd RepeatDays) Apply(start time.Time) time.Time {
	return start.AddDate(0, 0, rd.days)
}

func (rd *RepeatDays) Validate(repeat string) error {

	// значение repeat гарантированно не пустое
	parts := strings.FieldsFunc(repeat, func(c rune) bool {
		return unicode.IsSpace(c) || (c == ',')
	})

	// parts[0] гарантированно d, игнорируется
	params := parts[1:]

	// ключ не содержит параметров
	if len(params) == 0 {
		return ErrRuleEmptyKey
	}

	// проверка предельных значений ключа
	days, err := strconv.Atoi(params[0])
	if err != nil {
		return ErrRuleInvalidKey
	}
	if days < minDays || days > maxDays {
		return ErrRuleInvalidKey
	}

	rd.days = days

	return nil
}
