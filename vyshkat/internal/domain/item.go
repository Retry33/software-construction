package domain

var (
	_ InventoryItem  = (*Helmet)(nil)
	_ NumberAssigner = (*Helmet)(nil)

	_ InventoryItem  = (*DockStation)(nil)
	_ NumberAssigner = (*DockStation)(nil)
	_ EnergyConsumer = (*DockStation)(nil)

	_ InventoryItem  = (*ChargingCabinet)(nil)
	_ NumberAssigner = (*ChargingCabinet)(nil)
	_ EnergyConsumer = (*ChargingCabinet)(nil)
)

type Helmet struct {
	inventory
}

func NewHelmet(name string) (*Helmet, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	return &Helmet{inventory: inv}, nil
}

type DockStation struct {
	inventory
	energyUsage
}

func NewDockStation(name string, dailyKWh float64) (*DockStation, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return nil, err
	}
	return &DockStation{inventory: inv, energyUsage: energy}, nil
}

type ChargingCabinet struct {
	inventory
	energyUsage
}

func NewChargingCabinet(name string, dailyKWh float64) (*ChargingCabinet, error) {
	inv, err := newInventory(name)
	if err != nil {
		return nil, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return nil, err
	}
	return &ChargingCabinet{inventory: inv, energyUsage: energy}, nil
}
