package rules

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewRepeatMonths() *RepeatMonths {
	return &RepeatMonths{}
}

type RepeatMonths struct {
	days   []int
	months []int
}

func (rd *RepeatMonths) Validate(params []string) error {

	days0 := strings.FieldsFunc(params[0], func(c rune) bool {
		return (c == ',')
	})

	for _, v := range days0 {

		d, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("некорректное значение %s", v)
		}
		if d > 31 {
			return fmt.Errorf("некорректное значение %s", v)
		}

		if d < -2 {
			return fmt.Errorf("некорректное значение %s", v)
		}

		rd.days = append(rd.days, d)
	}

	sort.Ints(rd.days)

	if len(params) > 1 {

		months0 := strings.FieldsFunc(params[1], func(c rune) bool {
			return (c == ',')
		})

		for _, v := range months0 {

			m, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("некорректное значение %s", v)
			}
			if m < 1 || m > 12 {
				return fmt.Errorf("некорректное значение %s", v)
			}
			rd.months = append(rd.months, m)

		}
		sort.Ints(rd.months)
	}

	return nil
}

// stackoverflow
func DaysInMonth(t time.Time) int {
	y, m, _ := t.Date()
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (rd RepeatMonths) Apply(start time.Time) time.Time {

	var (
		nextYear, nextMonth, nextDay int
	)
	nextYear = start.Year()
	nextMonth = int(start.Month())
	nextDay = start.Day()

	daysInMonth := DaysInMonth(start)

	if len(rd.months) > 0 {
		for _, month := range rd.months {
			if month > int(start.Month()) {
				nextMonth = month
				goto todays
			}
		}
		if rd.months[0] <= int(start.Month()) {
			nextYear++
			nextMonth = rd.months[0]
		}
	}

todays:

	for _, day := range rd.days {
		if day > start.Day() {
			if day <= daysInMonth {
				nextDay = day
				goto finish
			}
			nextMonth++
			nextDay = day
			goto finish
		}
	}

	if rd.days[0] < 0 {
		nextDay = daysInMonth + rd.days[0] + 1 //нормализация -1, -2
		goto finish
	}

	if rd.days[0] <= start.Day() {
		if len(rd.months) == 0 {
			nextMonth++
		}
		nextDay = rd.days[0]
		goto finish
	}

finish:
	return time.Date(nextYear, time.Month(nextMonth), nextDay,
		start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), start.Location())
}
