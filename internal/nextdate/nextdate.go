package nextdate

import (
	"errors"
	"time"

	"github.com/dgshulgin/go_final_project/internal/nextdate/rules"
)

const (
	formatDateTime = "20060102"
)

var (
	ErrNextDateBeforeNow       = errors.New("дата меньше сегодняшней")
	ErrNextDateParsing         = errors.New("ошибка преобразования даты")
	ErrNextDateIncorrectRepeat = errors.New("некорректное значение поля repeat")
)

var ruler = map[string]rules.RepeatRule{
	"d": rules.NewRepeatDays(),
	"y": rules.NewRepeatYears(),
	"w": rules.NewRepeatWeeks(),
	"m": rules.NewRepeatMonths(),
}

// Validate проверяет поступившие значения date, now и вызывает RepeatRule.validate для проверки значения repeat.
func Validate(start string, now string, repeat string) error {

	if len(repeat) > 0 {
		// название ключа в списке допустимых: d, y, w, m
		rule, ok := ruler[string(repeat[0])]
		if !ok {
			return ErrNextDateIncorrectRepeat
		}

		err := rule.Validate(repeat)
		if err != nil {
			return ErrNextDateIncorrectRepeat
		}
	}

	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return ErrNextDateParsing
	}

	now0, err := time.Parse(formatDateTime, now)
	if err != nil {
		return ErrNextDateParsing
	}

	if start0.Before(now0) {
		return ErrNextDateBeforeNow
	}

	return nil
}

// NextDate вычисляет следующую дату события.
// Функция расчитывает что значение repeat уже обработано во время вызова Validate и состояние экземпляра RepeatRule установлено. Дополнительные проверки входных данных не выполняются.
func NextDate(start string, now string, repeat string) (string, error) {

	if len(repeat) == 0 {
		return time.Now().Format(formatDateTime), nil
	}

	// валидность уже проверили, игнорируем ошибки
	start0, _ := time.Parse(formatDateTime, start)
	now0, _ := time.Parse(formatDateTime, now)

	// валидность уже проверили, ключ присутствует
	rule, _ := ruler[string(repeat[0])]

	var nextTime time.Time
	nextTime = start0

	for {
		nextTime = rule.Apply(nextTime)
		if nextTime.After(now0) {
			break
		}
	}
	rule.Reset()

	return nextTime.Format(formatDateTime), nil
}
