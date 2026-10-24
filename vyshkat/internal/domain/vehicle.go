package domain

var (
	_ Vehicle        = (*Scooter)(nil)
	_ EnergyConsumer = (*Scooter)(nil)
	_ InventoryItem  = (*Scooter)(nil)
	_ NumberAssigner = (*Scooter)(nil)

	_ Vehicle        = (*Bicycle)(nil)
	_ InventoryItem  = (*Bicycle)(nil)
	_ NumberAssigner = (*Bicycle)(nil)

	_ Vehicle        = (*ElectricBike)(nil)
	_ EnergyConsumer = (*ElectricBike)(nil)
	_ InventoryItem  = (*ElectricBike)(nil)
	_ NumberAssigner = (*ElectricBike)(nil)
)

type Scooter struct {
	inventory
	energyUsage
	rental
}

func NewScooter(name string, dailyKWh float64, simplicity Simplicity) (*Scooter, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return nil, err
	}
	r, err := newRental(simplicity)
	if err != nil {
		return nil, err
	}
	return &Scooter{inventory: inv, energyUsage: energy, rental: r}, nil
}

type Bicycle struct {
	inventory
	rental
}

func NewBicycle(name string, simplicity Simplicity) (*Bicycle, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	r, err := newRental(simplicity)
	if err != nil {
		return nil, err
	}
	return &Bicycle{inventory: inv, rental: r}, nil
}

type ElectricBike struct {
	inventory
	energyUsage
	rental
}

func NewElectricBike(name string, dailyKWh float64, simplicity Simplicity) (*ElectricBike, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return nil, err
	}
	r, err := newRental(simplicity)
	if err != nil {
		return nil, err
	}
	return &ElectricBike{inventory: inv, energyUsage: energy, rental: r}, nil
}
