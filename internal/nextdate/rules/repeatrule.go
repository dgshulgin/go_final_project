package rules

import "time"

type RepeatRule interface {
	Validate(params []string) error
	Apply(start time.Time) time.Time
	Reset()
}
