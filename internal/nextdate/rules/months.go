package rules

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Реализация интерфейса RepeatRule для правила повторения m.
// Вычисляет ближайшую дату события, если условие повторения задано
// как дни месяца/месяцы.
// Правила формирования условия repeat:
// Формат: "m <дни> <месяцы>"
// <дни> - задача назначается на указанные дни месяца. Дни месяца
// перечислены через запятую. Допустимо указание значений для -1 и -2 для указания на последнее или предпоследнее число месяца, соответственно.
// Например, m 14 или m 1, 15, 25 или m -1
// <месяцы> - задача указывает на определенные номера месяцев. Номера месяцев перечисляются через запятую.
// Например, m 1,-1 2,8
type RepeatMonths struct {
	days   []int
	months []int
}

func NewRepeatMonths() *RepeatMonths {
	return &RepeatMonths{}
}

// Применение данной функции обязательно
func (rd *RepeatMonths) Reset() {
	rd.days = nil
	rd.months = nil
}

// Validate проверяет корректность значения repeat и заполняет
// слайсы rd.days и rd.months списками соотв параметров, которые обрабатываются
// уже на этапе Apply.
// Формат условия repeat: "m <дни> <месяцы>"
func (rd *RepeatMonths) Validate(repeat string) error {

	repeat0 := strings.ReplaceAll(repeat, " ", "#")
	parts := strings.FieldsFunc(repeat0, func(c rune) bool { return (c == '#') })
	// для ключа m не заданы <дни> и/или <месяцы>
	if len(parts) < 2 {
		return ErrRuleEmptyKey
	}

	// Извлечь список дней как слайс
	days0 := strings.FieldsFunc(parts[1], func(c rune) bool {
		return (c == ',')
	})

	// Проверка корректности значений <дни>
	for _, v := range days0 {
		d, err := strconv.Atoi(v)
		if err != nil {
			return ErrRuleInvalidKey
		}
		if d > 31 || d < -2 {
			return ErrRuleInvalidKey
		}
		rd.days = append(rd.days, d)
	}
	sort.Ints(rd.days)

	// Параметр <месяцы> может отсутствовать
	if len(parts) == 3 {
		// Извлечь список месяцев как слайс
		months0 := strings.FieldsFunc(parts[2], func(c rune) bool {
			return (c == ',')
		})
		// Проверка корректности значений <месяцы>
		for _, v := range months0 {
			m, err := strconv.Atoi(v)
			if err != nil {
				return ErrRuleInvalidKey
			}
			if m < 1 || m > 12 {
				return ErrRuleInvalidKey
			}
			rd.months = append(rd.months, m)
		}
		sort.Ints(rd.months)
	}
	return nil
}

// stackoverflow
func daysInMonth(t time.Time) int {
	y, m, _ := t.Date()
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func filter(elems []int, test func(val int) bool) ([]int, bool) {
	r0 := []int{}
	for _, val := range elems {
		if test(val) {
			r0 = append(r0, val)
		}
	}
	return r0, len(r0) == 0
}

// Apply однократно вычисляет следующую дату события исходя из значений start, rd.days, rd.months. Функция NextDate итеративно вызывает Apply, таким образом происходит пошаговое "приближение" к фактической следующей дате.
func (rd RepeatMonths) Apply(start time.Time) time.Time {

	nextYear := start.Year()
	nextMonth := int(start.Month())
	nextDay := start.Day()

	daysInMonth := daysInMonth(start)

	// заменить отрицательные значения -1 и -2 на фактическую дату последнего/предпоследнего дня
	if rd.days[0] < 0 {
		for idx, val := range rd.days {
			if val < 0 {
				rd.days[idx] = daysInMonth + val + 1
			}
		}
		sort.Ints(rd.days)
	}

	if len(rd.months) > 0 {

		// в этом году, если есть месяцы больше текущего
		month0, empty := filter(rd.months, func(val int) bool {
			return val > nextMonth
		})
		if !empty {
			nextMonth = month0[0]
			nextDay = rd.days[0]
			return time.Date(nextYear, time.Month(nextMonth), nextDay,
				start.Hour(), start.Minute(), start.Second(),
				start.Nanosecond(),
				start.Location())
		}
		// уже в след году
		nextYear++
		nextMonth = rd.months[0]
		nextDay = rd.days[0]
		return time.Date(nextYear, time.Month(nextMonth), nextDay,
			start.Hour(), start.Minute(), start.Second(),
			start.Nanosecond(),
			start.Location())
	}

	// в этом месяце, если есть дни больше текущего
	days0, empty := filter(rd.days, func(val int) bool {
		return val <= daysInMonth && val > nextDay
	})
	if !empty {
		nextDay = days0[0]
		return time.Date(nextYear, time.Month(nextMonth), nextDay,
			start.Hour(), start.Minute(), start.Second(),
			start.Nanosecond(),
			start.Location())
	}
	// уже в следующем месяце
	nextMonth++
	nextDay = rd.days[0]
	return time.Date(nextYear, time.Month(nextMonth), nextDay,
		start.Hour(), start.Minute(), start.Second(),
		start.Nanosecond(),
		start.Location())
}
