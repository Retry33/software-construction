package inspection

import (
	"errors"
	"fmt"

	"vyshkat/internal/domain"
)

var (
	ErrNilVehicle           = errors.New("no vehicle presented for inspection")
	ErrExcessiveConsumption = errors.New("daily energy consumption exceeds the allowed limit")
)

type ServiceCenter struct {
	maxDailyKWh float64
}

func NewServiceCenter(maxDailyKWh float64) *ServiceCenter {
	return &ServiceCenter{maxDailyKWh: maxDailyKWh}
}

func (s *ServiceCenter) Inspect(v domain.Vehicle) error {
	if v == nil {
		return ErrNilVehicle
	}
	if c, ok := v.(domain.EnergyConsumer); ok && c.DailyEnergyKWh() > s.maxDailyKWh {
		return fmt.Errorf("%w: %.2f kWh > %.2f kWh", ErrExcessiveConsumption, c.DailyEnergyKWh(), s.maxDailyKWh)
	}
	return nil
}
