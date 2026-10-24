package console

import "vyshkat/internal/domain"

type vehicleKind struct {
	label      string
	usesEnergy bool
	build      func(name string, dailyKWh float64, s domain.Simplicity) (domain.Vehicle, error)
}

type itemKind struct {
	label      string
	usesEnergy bool
	build      func(name string, dailyKWh float64) (domain.InventoryItem, error)
}

var vehicleKinds = []vehicleKind{
	{
		label:      "Electric scooter",
		usesEnergy: true,
		build: func(name string, kwh float64, s domain.Simplicity) (domain.Vehicle, error) {
			return asVehicle(domain.NewScooter(name, kwh, s))
		},
	},
	{
		label: "Bicycle",
		build: func(name string, _ float64, s domain.Simplicity) (domain.Vehicle, error) {
			return asVehicle(domain.NewBicycle(name, s))
		},
	},
	{
		label:      "Electric bike",
		usesEnergy: true,
		build: func(name string, kwh float64, s domain.Simplicity) (domain.Vehicle, error) {
			return asVehicle(domain.NewElectricBike(name, kwh, s))
		},
	},
}

var itemKinds = []itemKind{
	{
		label: "Helmet",
		build: func(name string, _ float64) (domain.InventoryItem, error) {
			return asItem(domain.NewHelmet(name))
		},
	},
	{
		label:      "Dock station",
		usesEnergy: true,
		build: func(name string, kwh float64) (domain.InventoryItem, error) {
			return asItem(domain.NewDockStation(name, kwh))
		},
	},
	{
		label:      "Charging cabinet",
		usesEnergy: true,
		build: func(name string, kwh float64) (domain.InventoryItem, error) {
			return asItem(domain.NewChargingCabinet(name, kwh))
		},
	},
}

func asVehicle[T domain.Vehicle](v T, err error) (domain.Vehicle, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

func asItem[T domain.InventoryItem](item T, err error) (domain.InventoryItem, error) {
	if err != nil {
		return nil, err
	}
	return item, nil
}
