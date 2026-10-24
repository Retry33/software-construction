package fleet

import "vyshkat/internal/domain"

type Inspector interface {
	Inspect(v domain.Vehicle) error
}
