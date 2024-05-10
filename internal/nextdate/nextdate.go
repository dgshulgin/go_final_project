package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	minDays        = 1
	maxDays        = 400
	stepYears      = 1
	weekMonday     = 1
	weekSunday     = 7
	formatDateTime = "20060102"
)

// Тип функции-обработчика для ключей d, w, m, y и пустого значения repeat ("")
type Handler func(start string, params []string) (string, error)

var ruler = map[string]Handler{
	"d": repeatInDays,
	"y": repeatInYears,
	"w": repeatInWeeks,
	"m": repeatInMonths,
	"":  repeatIsEmpty,
}

// Вычисляет ближайщую дату наступления события.
// start - текущая дата наступления события
// now - сегодняшная (на момент вызова функции) дата
// repeat - условия назначения ближайшей даты наступления события
func NextDate(start string, now string, repeat string) (string, error) {

	//текущая дата не задана, ближайшей датой считается сегодняшняя дата
	if len(start) == 0 {
		return now, nil
	}

	// условия повторения не определены
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

	// допустимые ключи повторения: d, y, w, m
	handler, ok := ruler[rep[0]]
	if !ok {
		return "", fmt.Errorf("некорректное значение поля repeat: %s", rep[0])
	}

	// параметры ключа не заданы. Допустимо для ключа y, для остальных - считается ошибкой.
	repTail := rep[1:]
	if len(repTail) == 0 {
		repTail = []string{}
	}

	// Вспомогательная функция NextDate0 циклически вычисляет ближайшую дату
	// используя выбранную функцию-обработчик (handler)
	return NextDate0(start, now, repTail, handler)
}

// Циклически вычисляет ближайшую дату наступления события с помощью функции-обработчика handler.
func NextDate0(
	start string,
	now string,
	params []string,
	handler func(start0 string, params0 []string) (string, error)) (string, error) {

	var nextTime string
	var err error
	nextTime = start
	for {
		nextTime, err = handler(nextTime, params)
		if err != nil {
			return "", fmt.Errorf("ошибка вычисления даты, %s", err.Error())
		}
		nt0, _ := time.Parse(formatDateTime, nextTime)
		now0, _ := time.Parse(formatDateTime, now)
		if nt0.After(now0) {
			break
		}
	}
	return nextTime, nil
}

// Вычисляет ближайшую дату события, если условие повторения задано как кол-во дней.
// Правила формирования условия repeat:
// Формат: "d <число>", где число находится в пределах 1-400
func repeatInDays(start string, params []string) (string, error) {

	// проверка формата
	if len(params) == 0 {
		return "", fmt.Errorf("отсутствует значение для ключа d")
	}

	// проверка предельных значений ключа
	days, err := strconv.Atoi(params[0])
	if err != nil {
		return "", fmt.Errorf("некорректное значение аргумента для ключа d: %s", params[0])
	}
	if days < minDays || days > maxDays {
		return "", fmt.Errorf("недопустимое количество дней для ключа d %d (%d-%d))", days, minDays, maxDays)
	}

	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты %s", start)
	}

	nextTime := start0.AddDate(0, 0, days)

	return nextTime.Format(formatDateTime), nil
}

// Вычисляет ближайшую дату события, если задача выполняется ежегодно.
// Выполнение задачи переносится на один год вперед.
// Правила формирования условия repeat:
// Формат: "y"
func repeatInYears(start string, params []string) (string, error) {

	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты %s", start)
	}

	nextTime := start0.AddDate(stepYears, 0, 0)

	return nextTime.Format(formatDateTime), nil
}

// Вычисляет ближайщую дату события, если задача назначена в указанные дни недели, ключ w.
// Формат: "w <число>,<число>", где 1-понедельник, 7-воскресенье
// Пример:
func repeatInWeeks(start string, params []string) (string, error) {
	return "", fmt.Errorf("неподдерживаемый формат w")

	// // проверка формата
	// if len(params) == 0 {
	// 	return "", fmt.Errorf("отсутствует значение для ключа w")
	// }
	// // проверка предельных значений
	// for _, par := range params {

	// 	day, err := strconv.Atoi(par)
	// 	if err != nil {
	// 		return "", fmt.Errorf("некорректное значение аргумента для ключа w %d", day)
	// 	}

	// 	if day < weekMonday || day > weekSunday {
	// 		return "", fmt.Errorf("некорректное значение аргумента для ключа w %d", day)
	// 	}
	// }

	// // какой день недели соотв start
	// start0, err := time.Parse(formatDateTime, start)
	// if err != nil {
	// 	return "", fmt.Errorf("ошибка преобразования даты %s", start)
	// }
	// startWeekDay := start0.Weekday()
	// if startWeekDay == time.Sunday {
	// 	startWeekDay = weekSunday
	// }
	// fmt.Printf("%s is weekday %d\n", start, startWeekDay)

	// sort.Strings(params)

	// for _, par := range params {
	// 	day, _ := strconv.Atoi(par)

	// 	if day > int(startWeekDay) {
	// 		return start0.AddDate(0, 0, day-int(startWeekDay)).Format(formatDateTime), nil
	// 	}
	// }
	// dd, _ := strconv.Atoi(params[0])
	// if dd == int(startWeekDay) {
	// 	return start0.AddDate(0, 0, 6).Format(formatDateTime), nil
	// }
	// count := 7 - int(startWeekDay) + dd
	// return start0.AddDate(0, 0, count).Format(formatDateTime), nil

	// //задача перносится на ближайший день недели, который наступит раньше

	// // fmt.Printf("count=%d\n", count)

	// // nextTime := start0.AddDate(0, 0, count)

	// // fmt.Printf("nextTime is %s\n", nextTime.Format(formatDateTime))

	// // return nextTime.Format(formatDateTime), nil

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
