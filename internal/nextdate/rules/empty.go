package rules

import "time"

// Правило обработки пустого ключа (значение repeat не задано)
type RepeatEmpty struct {
}

func NewRepeatEmpty() *RepeatEmpty { return &RepeatEmpty{} }

func (rd RepeatEmpty) Validate(params []string) error {
	return nil
}
func (rd RepeatEmpty) Apply(start time.Time) time.Time {

	if start.Before(time.Now()) {
		return time.Now()
	}
	return start
}
