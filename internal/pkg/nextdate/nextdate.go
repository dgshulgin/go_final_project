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
	if len(params) == 0 {
		return "", fmt.Errorf("отсутствует значение для ключа d")
	}
	days, err := strconv.Atoi(params[0])
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

var ruler = map[string]Handler{
	"d": repeatInDays,
	"y": repeatInYears,
	"w": repeatInWeeks,
	"m": repeatInMonths,
}

// Ожидает непустой repeat.
func NextDate(start string, now string, repeat string) (string, error) {
	rep := strings.FieldsFunc(repeat, func(c rune) bool {
		return unicode.IsSpace(c) || (c == ',')
	})
	handler, ok := ruler[rep[0]]
	if !ok {
		return "", fmt.Errorf("некорректное значение поля repeat: %s", rep[0])
	}
	rep = rep[1:]
	if len(rep) == 0 {
		rep = []string{}
	}
	nextTime, err := handler(start, rep)
	if err != nil {
		return "", fmt.Errorf("ошибка вычисления даты, %w", err)
	}

	nt0, _ := time.Parse(formatDateTime, nextTime)
	now0, _ := time.Parse(formatDateTime, now)
	if nt0.After(now0) {
		return nextTime, nil
	}
	return NextDate(nextTime, now, repeat)
}
