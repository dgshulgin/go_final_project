package rules

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	weekMonday = 1
	weekSunday = 7
)

// Реализация интерфейса RepeatRule для правила повторения w.
// Вычисляет ближайщую дату события, если задача назначена в указанные дни недели, ключ w.
// Формат: "w <число>,<число>", где 1-понедельник, 7-воскресенье
type RepeatWeeks struct {
	weekDays []int
}

func NewRepeatWeeks() *RepeatWeeks { return &RepeatWeeks{} }

// Применение данной функции обязательно
func (rd *RepeatWeeks) Reset() {
	rd.weekDays = nil
}

// Validate проверяет корректность параметров ключа w и заполняет
// слайс rd.weekDays значениями, которые обрабатываются
// уже на этапе Apply.
// Формат: "w <число>,<число>", где 1-понедельник, 7-воскресенье
func (rd *RepeatWeeks) Validate(repeat string) error {

	parts := strings.FieldsFunc(repeat, func(c rune) bool {
		return unicode.IsSpace(c) || (c == ',')
	})

	// parts[0] гарантированно w, игнорируется
	params := parts[1:]
	// ключ не содержит параметров
	if len(params) == 0 {
		return ErrRuleEmptyKey
	}

	// проверка предельных значений ключа
	for _, par := range params {

		day, err := strconv.Atoi(par)
		if err != nil {
			return ErrRuleInvalidKey
		}

		if day < weekMonday || day > weekSunday {
			return ErrRuleInvalidKey
		}

		rd.weekDays = append(rd.weekDays, day)
	}
	sort.Ints(rd.weekDays)

	return nil
}

// Порядковые номера дней недели в rd.weekDays должны быть отсортированы по возрастанию. В таком случае, первое значение, превышающее start.Weekday является ближайшим днем для наступления события на этой неделе.
// Если таковое не найдено, то первый элемент слайса rd.WeekDays (минимальное значение) является ближайшим днем для наступления события на следующей неделе.
// Еще возможен вариант "w 3" при startWeekDay == 3, означает, что событие произойдет точно через неделю.
// Оригинальный time.Weekday принимает Воскресенье за 0, тогда как в задании порядковый номер для Воскресенья равен 7. Преобразование в первых строках функции.
func (rd RepeatWeeks) Apply(start time.Time) time.Time {

	startWeekDay := start.Weekday()
	if startWeekDay == time.Sunday {
		startWeekDay = weekSunday
	}

	for _, day := range rd.weekDays {
		if day > int(startWeekDay) {
			return start.AddDate(0, 0, day-int(startWeekDay))
		}
	}
	first := rd.weekDays[0]
	if first == int(startWeekDay) {
		return start.AddDate(0, 0, 6)
	}
	count := 7 - int(startWeekDay) + first
	return start.AddDate(0, 0, count)

}
