package rules

import (
	"fmt"
	"time"
)

func NewRepeatMonths() *RepeatMonths { return &RepeatMonths{} }

type RepeatMonths struct {
}

func (rd RepeatMonths) Validate(params []string) error {
	return fmt.Errorf("неподдерживаемый формат m")
}
func (rd RepeatMonths) Apply(start time.Time) time.Time {
	return time.Now()
}
