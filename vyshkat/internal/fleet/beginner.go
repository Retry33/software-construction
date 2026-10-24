package fleet

import "vyshkat/internal/domain"

const DefaultBeginnerMinSimplicity = 6

type BeginnerPolicy interface {
	Allows(v domain.Vehicle) bool
}

type MinSimplicityPolicy struct {
	Min int
}

func (p MinSimplicityPolicy) Allows(v domain.Vehicle) bool {
	return v.Simplicity().Value() >= p.Min
}

func DefaultBeginnerPolicy() BeginnerPolicy {
	return MinSimplicityPolicy{Min: DefaultBeginnerMinSimplicity}
}
