package domain

import "fmt"

const (
	MinSimplicity = 1
	MaxSimplicity = 10
)

type Simplicity struct {
	value int
}

func NewSimplicity(value int) (Simplicity, error) {
	if value < MinSimplicity || value > MaxSimplicity {
		return Simplicity{}, fmt.Errorf("%w: got %d", ErrInvalidSimplicity, value)
	}
	return Simplicity{value: value}, nil
}

func (s Simplicity) Value() int {
	return s.value
}

func (s Simplicity) valid() bool {
	return s.value >= MinSimplicity && s.value <= MaxSimplicity
}
