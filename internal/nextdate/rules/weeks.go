package rules

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

const (
	weekMonday = 1
	weekSunday = 7
)

// Вычисляет ближайщую дату события, если задача назначена в указанные дни недели, ключ w.
// Формат: "w <число>,<число>", где 1-понедельник, 7-воскресенье
// Пример:
// задача перносится на ближайший день недели, который наступит раньше

type RepeatWeeks struct {
	weekDays []int
}

func NewRepeatWeeks() *RepeatWeeks { return &RepeatWeeks{} }

func (rd *RepeatWeeks) Validate(params []string) error {

	// проверка формата
	if len(params) == 0 {
		return fmt.Errorf("отсутствует значение для ключа w")
	}
	// проверка предельных значений
	for _, par := range params {

		day, err := strconv.Atoi(par)
		if err != nil {
			return fmt.Errorf("некорректное значение аргумента для ключа w %d", day)
		}

		if day < weekMonday || day > weekSunday {
			return fmt.Errorf("некорректное значение аргумента для ключа w %d", day)
		}

		rd.weekDays = append(rd.weekDays, day)
	}

	// Apply ожидает что порядковые номера дней недели расположены по возрастанию
	sort.Ints(rd.weekDays)

	return nil
}

// Порядковые номера дней недели в rd.weekDays должны быть отсортированы по возрастанию.
// В таком случае, первое значение, превышающее start.Weekday является ближайшим днем
// для наступления события на этой неделе.
// Если таковое не найдено, то первый элемент слайса rd.WeekDays (минимальное значение)
// является ближайшим днем для наступления события на следующей неделе.
// Еще возможен вариант "w 3" при startWeekDay == 3, означает, что событие произойдет
// точно через неделю.
// Оригинальный time.Weekday принимает Воскресенье за 0, тогда как в задании порядковый
// номер для Воскресенья равен 7. Преобразование в первых строках функции.
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
