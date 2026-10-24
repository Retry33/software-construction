package domain

import (
	"math"
	"strings"
)

type inventory struct {
	name   string
	number string
}

func newInventory(name string) (inventory, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return inventory{}, ErrEmptyName
	}
	return inventory{name: name}, nil
}

func (i *inventory) Name() string { return i.name }

func (i *inventory) Number() string { return i.number }

func (i *inventory) AssignNumber(number string) error {
	if strings.TrimSpace(number) == "" {
		return ErrEmptyNumber
	}
	if i.number != "" {
		return ErrNumberAlreadyAssigned
	}
	i.number = number
	return nil
}

type energyUsage struct {
	dailyKWh float64
}

func newEnergyUsage(dailyKWh float64) (energyUsage, error) {
	if dailyKWh < 0 || math.IsNaN(dailyKWh) || math.IsInf(dailyKWh, 0) {
		return energyUsage{}, ErrNegativeEnergy
	}
	return energyUsage{dailyKWh: dailyKWh}, nil
}

func (e *energyUsage) DailyEnergyKWh() float64 { return e.dailyKWh }

type rental struct {
	simplicity Simplicity
}

func newRental(s Simplicity) (rental, error) {
	if !s.valid() {
		return rental{}, ErrInvalidSimplicity
	}
	return rental{simplicity: s}, nil
}

func (r *rental) Simplicity() Simplicity { return r.simplicity }
