package nextdate

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/dgshulgin/go_final_project/internal/nextdate/rules"
)

const (
	formatDateTime = "20060102"
)

var ruler = map[string]rules.RepeatRule{}

// Вычисляет ближайщую дату наступления события.
// start - текущая дата наступления события
// now - сегодняшная (на момент вызова функции) дата
// repeat - правила назначения ближайшей даты наступления события
func NextDate(start string, now string, repeat string) (string, error) {

	ruler["d"] = rules.NewRepeatDays()
	ruler["y"] = rules.NewRepeatYears()
	ruler["w"] = rules.NewRepeatWeeks()
	ruler["m"] = rules.NewRepeatMonths()
	ruler[""] = rules.NewRepeatEmpty()

	//текущая дата не задана, ближайшей датой считается сегодняшняя дата
	if len(start) == 0 {
		return now, nil
	}

	start0, err := time.Parse(formatDateTime, start)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты %s", start)
	}

	now0, err := time.Parse(formatDateTime, now)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования сегодняшней даты %s", now)
	}

	// условия повторения не определены (пустой ключ)
	if len(repeat) == 0 {
		err := ruler[""].Validate([]string{})
		if err != nil {
			return "", err
		}
		nextTime := ruler[""].Apply(start0)
		return nextTime.Format(formatDateTime), nil
	}

	rep := strings.FieldsFunc(repeat, func(c rune) bool {
		return unicode.IsSpace(c) || (c == ',')
	})

	// допустимые ключи повторения: d, y, w, m
	repeatRule, ok := ruler[rep[0]]
	if !ok {
		return "", fmt.Errorf("некорректное значение поля repeat: %s", rep[0])
	}

	// параметры ключа не заданы. Допустимо для ключа y, для остальных - считается ошибкой.
	repTail := rep[1:]
	if len(repTail) == 0 {
		repTail = []string{}
	}

	// ключ m  устроен сложнее других: ключ<sp>дни<sp>месяцы разделены пробелом
	// выполним замену символа пробела на символ # (в формате никогда не встречается)
	// получившуюся строку разделим на ключ и параметры: "m" "дни" "месяцы"
	// итоговое наполнение списка ключей в repTail отличается от других ключей
	// параметры должны быть дополнительно обработаны в валидаторе для ключа m
	// допущение: в строке не содержится пробелов за искл разделителей ключ-дни-месяцы
	if strings.Compare(rep[0], "m") == 0 {

		repeat0 := strings.ReplaceAll(repeat, " ", "#")
		parts := strings.FieldsFunc(repeat0, func(c rune) bool { return (c == '#') })

		repTail = nil // сбросить предыдущее содержимое слайса
		repTail = append(repTail, parts[1])
		//а вот месяцев может и не быть
		if len(parts) > 2 {
			repTail = append(repTail, parts[2])
		}
	}

	// Вспомогательная функция NextDate0 циклически вычисляет ближайшую дату
	// используя выбранное правило repeatRule
	nd, err := NextDate0(start0, now0, repTail, repeatRule)
	if err != nil {
		return "", err
	}

	return nd.Format(formatDateTime), nil
}

// Циклически вычисляет ближайшую дату наступления события с помощью функции-обработчика handler.
func NextDate0(
	start time.Time,
	now time.Time,
	params []string,
	rule rules.RepeatRule) (time.Time, error) {

	err := rule.Validate(params)
	if err != nil {
		return time.Now(), fmt.Errorf("некорректное значение поля repeat, %s", err.Error())
	}

	var nextTime time.Time
	nextTime = start
	for {
		nextTime = rule.Apply(nextTime)
		if nextTime.After(now) {
			break
		}
	}
	return nextTime, nil
}
