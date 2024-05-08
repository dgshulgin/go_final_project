package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	maxDays        = 400
	formatDateTime = "20060102"
)

type Handler func(start string, params []string) (string, error)

// допустимо только одно значение < 400
func repeatInDays(start string, params []string) (string, error) {
	///fmt.Printf("repeatInDays, start=%v, params=%v\n", start, params)
	if len(params) == 0 {
		return "", fmt.Errorf("отсутствует значение для ключа d")
	}
	days, err := strconv.Atoi(params[0])
	//fmt.Printf("repeatInDays, days=%v\n", days)
	if err != nil {
		return "", fmt.Errorf("некорректное значение аргумента для ключа d: %s", params[0])
	}
	if max(days, maxDays) == days {
		return "", fmt.Errorf("значение аргумента для ключа d превышает 400")
	}
	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты %s", start)
	}
	nextTime := start0.AddDate(0, 0, days)
	//fmt.Printf("repeatInDays, nextTime=%v\n", nextTime)
	return nextTime.Format(formatDateTime), nil
}

func repeatInYears(start string, params []string) (string, error) {
	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты %s", start)
	}
	nextTime := start0.AddDate(1, 0, 0)
	return nextTime.Format(formatDateTime), nil

}

func repeatInWeeks(start string, params []string) (string, error) {
	return "", fmt.Errorf("неподдерживаемый формат w")
}

func repeatInMonths(start string, params []string) (string, error) {
	return "", fmt.Errorf("неподдерживаемый формат m")
}

func repeatIsEmpty(start string, params []string) (string, error) {
	//fmt.Printf("repeatIsEmpty, start = %v\n", start)
	startDate, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("некорректное значение поля date") //, %v", err)
	}
	//fmt.Printf("repeatIsEmpty, parsed start=%v\n", startDate)
	nowDate := time.Now()
	if startDate.Before(nowDate) {
		return nowDate.Format("20060102"), nil
	}
	return start, nil
}

var ruler = map[string]Handler{
	"d": repeatInDays,
	"y": repeatInYears,
	"w": repeatInWeeks,
	"m": repeatInMonths,
	"":  repeatIsEmpty,
}

/*
3. date >= now, repeat == empty => date
date >= now, repeat => nextdate
2. date < now, repeat is empty => now
date < now, repeat => nextdate
1. date is empty, repeat is empty => now
1. date is empty, repeat => now
*/
func NextDate0(start string, now string, params []string, handler func(start0 string, params0 []string) (string, error)) (string, error) {
	// nextTime, err := handler(start, params)
	// if err != nil {
	// 	return "", fmt.Errorf("ошибка вычисления даты, %w", err)
	// }
	// nt0, _ := time.Parse(formatDateTime, nextTime)
	// now0, _ := time.Parse(formatDateTime, now)
	// if nt0.After(now0) {
	// 	return nextTime, nil
	// }
	// return NextDate0(nextTime, now, params, handler)
	var nextTime string
	var err error
	nextTime = start
	for {
		nextTime, err = handler(nextTime, params)
		if err != nil {
			return "", fmt.Errorf("ошибка вычисления даты, %w", err)
		}
		nt0, _ := time.Parse(formatDateTime, nextTime)
		now0, _ := time.Parse(formatDateTime, now)
		if nt0.After(now0) {
			break
		}
	}
	return nextTime, nil
}

func NextDate(start string, now string, repeat string) (string, error) {
	if len(start) == 0 {
		return now, nil
		//n0, _ := time.Parse("20060102", now)

		//return n0.AddDate(0, 0, 1).Format("20060102"), nil
	}
	// startDate, err := time.Parse(formatDateTime, start)
	// if err != nil {
	// 	return "", fmt.Errorf("некорректное значение поля date, %w", err)
	// }
	// nowDate, err := time.Parse(formatDateTime, now)
	// if err != nil {
	// 	return "", fmt.Errorf("некорректное значение аргумента now, %w", err)
	// }
	// if startDate.Before(nowDate) {
	// 	if len(repeat) == 0 {
	// 		return now, nil
	// 	}
	// }
	if len(repeat) == 0 {
		handler := ruler[""]
		nextTime, err := handler(start, nil)
		if err != nil {
			return "", err
		}
		return nextTime, nil
	}
	rep := strings.FieldsFunc(repeat, func(c rune) bool {
		return unicode.IsSpace(c) || (c == ',')
	})
	handler, ok := ruler[rep[0]]
	if !ok {
		return "", fmt.Errorf("некорректное значение поля repeat: %s", rep[0])
	}
	repTail := rep[1:]
	if len(repTail) == 0 {
		repTail = []string{}
	}
	return NextDate0(start, now, repTail, handler)

	// nextTime, err := handler(start, repTail)
	// if err != nil {
	// 	return "", fmt.Errorf("ошибка вычисления даты, %w", err)
	// }
	// nt0, _ := time.Parse(formatDateTime, nextTime)
	// now0, _ := time.Parse(formatDateTime, now)
	// if nt0.After(now0) {
	// 	return nextTime, nil
	// }
	// return NextDate(nextTime, now, repeat)
}

// func NextDate(start string, now string, repeat string) (string, error) {
// 	rep := strings.FieldsFunc(repeat, func(c rune) bool {
// 		return unicode.IsSpace(c) || (c == ',')
// 	})
// 	handler, ok := ruler[rep[0]]
// 	if !ok {
// 		return "", fmt.Errorf("некорректное значение поля repeat: %s", rep[0])
// 	}
// 	rep = rep[1:]
// 	if len(rep) == 0 {
// 		rep = []string{}
// 	}
// 	nextTime, err := handler(start, rep)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка вычисления даты, %w", err)
// 	}

// 	nt0, _ := time.Parse(formatDateTime, nextTime)
// 	now0, _ := time.Parse(formatDateTime, now)
// 	if nt0.After(now0) {
// 		return nextTime, nil
// 	}
// 	return NextDate(nextTime, now, repeat)
// }
