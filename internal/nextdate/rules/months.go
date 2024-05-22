package rules

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	formatDateTime = "20060102"
)

func NewRepeatMonths() *RepeatMonths {
	return &RepeatMonths{}
}

type RepeatMonths struct {
	days   []int
	months []int
}

func (rd *RepeatMonths) Reset() {
	rd.days = nil
	rd.months = nil
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
	fmt.Printf("Validate-1: rd.days=%v\n", rd.days)

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

	//start
	//now
	//nextDate := start
	nextYear := start.Year()
	nextMonth := int(start.Month())
	nextDay := start.Day()

	daysInMonth := DaysInMonth(start)

	filter := func(elems []int, test func(val int) bool) ([]int, bool) {
		r0 := []int{}
		for _, val := range elems {
			if test(val) {
				r0 = append(r0, val)
			}
		}
		return r0, len(r0) == 0
	}

	fmt.Printf("Validate-1.5: rd.days=%v\n", rd.days)

	//заменить отрицательные значения на реальные дни
	if rd.days[0] < 0 {
		for idx, val := range rd.days {
			if val < 0 {
				rd.days[idx] = daysInMonth + val + 1
			}
		}
		sort.Ints(rd.days)
	}

	fmt.Printf("Validate-2: rd.days=%v\n", rd.days)
	//nextYear = max(nextYear, rd.now.Year())
	if len(rd.months) > 0 {
		//lastMonth := rd.months[len(rd.months)-1]

		month0, empty := filter(rd.months, func(val int) bool {
			return val > nextMonth
		})
		// в этом году ?
		if !empty {
			nextMonth = month0[0]
			nextDay = rd.days[0]
			fmt.Printf("Apply-1: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)
			return time.Date(nextYear, time.Month(nextMonth), nextDay,
				start.Hour(), start.Minute(), start.Second(),
				start.Nanosecond(),
				start.Location())
		}
		//в след году
		nextYear++
		nextMonth = rd.months[0]
		nextDay = rd.days[0]
		fmt.Printf("Apply-2: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)
		return time.Date(nextYear, time.Month(nextMonth), nextDay,
			start.Hour(), start.Minute(), start.Second(),
			start.Nanosecond(),
			start.Location())
	}

	fmt.Printf("Validate-3: rd.days=%v\n", rd.days)
	//возможно в этом месяце ?
	//lastDay := rd.days[len(rd.days)-1]
	days0, empty := filter(rd.days, func(val int) bool {
		return val <= daysInMonth && val > nextDay
	})
	if !empty {
		nextDay = days0[0]
		fmt.Printf("Apply-3: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)

		return time.Date(nextYear, time.Month(nextMonth), nextDay,
			start.Hour(), start.Minute(), start.Second(),
			start.Nanosecond(),
			start.Location())
	}
	//только в следующем месяце
	nextMonth++
	nextDay = rd.days[0]
	fmt.Printf("Apply-4: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)

	return time.Date(nextYear, time.Month(nextMonth), nextDay,
		start.Hour(), start.Minute(), start.Second(),
		start.Nanosecond(),
		start.Location())
}

// if len(rd.months) > 0 {
// 	for nextYear < rd.now.Year() && nextMonth < int(rd.now.Month()) {
// 		lastMonth := rd.months[len(rd.months)-1]
// 		if lastMonth == nextMonth {
// 			nextYear++
// 			nextMonth = rd.months[0]
// 			return time.Date(nextYear, time.Month(nextMonth), nextDay,
// 				start.Hour(), start.Minute(), start.Second(),
// 				start.Nanosecond(),
// 				start.Location())
// 		}
// 		for _, month := range rd.months {
// 			if month > nextMonth {
// 				nextMonth = month
// 				return time.Date(nextYear, time.Month(nextMonth), nextDay,
// 					start.Hour(), start.Minute(), start.Second(),
// 					start.Nanosecond(),
// 					start.Location())
// 			}
// 		}
// 	}
// }

// //{"20231106", "m 13", "20240213"},
// for nextDay < rd.now.Day() {
// 	lastDay := rd.days[len(rd.days)-1]
// 	if lastDay == nextDay {
// 		nextMonth++
// 		nextDay = rd.days[0]
// 		return time.Date(nextYear, time.Month(nextMonth), nextDay,
// 			start.Hour(), start.Minute(), start.Second(),
// 			start.Nanosecond(),
// 			start.Location())
// 	}
// }

// for day := range rd.days {
// 	if day > nextDay {
// 		nextDay = day
// 		break
// 	}
// }

// //finish:
// return time.Date(nextYear, time.Month(nextMonth), nextDay,
// 	start.Hour(), start.Minute(), start.Second(),
// 	start.Nanosecond(),
// 	start.Location())

//}

// func (rd RepeatMonths) Apply(start time.Time) time.Time {

// 	var (
// 		nextYear, nextMonth, nextDay int
// 	)
// 	nextYear = start.Year()
// 	nextMonth = int(start.Month())
// 	nextDay = start.Day()

// 	daysInMonth := DaysInMonth(start)

// 	fmt.Printf("mApply: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)

// 	if len(rd.months) > 0 {
// 		for _, month := range rd.months {
// 			if month > int(start.Month()) {
// 				nextMonth = month
// 				goto todays
// 			}
// 		}
// 		if rd.months[0] <= int(start.Month()) {
// 			nextYear++
// 			nextMonth = rd.months[0]
// 		}
// 	}

// todays:
// 	fmt.Printf("mApply, after months: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)

// 	for _, day := range rd.days {
// 		if day > start.Day() {
// 			if day <= daysInMonth {
// 				nextDay = day
// 				goto finish
// 			}
// 			nextMonth++
// 			nextDay = day
// 			goto finish
// 		}
// 	}

// 	if rd.days[0] < 0 {
// 		nextDay = daysInMonth + rd.days[0] + 1 //нормализация -1, -2
// 		goto finish
// 	}

// 	if rd.days[0] <= start.Day() {
// 		if len(rd.months) == 0 {
// 			nextMonth++
// 		}
// 		nextDay = rd.days[0]
// 		goto finish
// 	}

// finish:
// 	fmt.Printf("mApply, after days: nextYear=%v, nextMonth=%v, nextDay=%v, daysInMonth=%v\n", nextYear, nextMonth, nextDay, daysInMonth)
// 	return time.Date(nextYear, time.Month(nextMonth), nextDay,
// 		start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), start.Location())
// }
