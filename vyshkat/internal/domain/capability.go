package domain

type InventoryItem interface {
	Name() string
	Number() string
}

type NumberAssigner interface {
	AssignNumber(number string) error
}

type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

type Vehicle interface {
	Simplicity() Simplicity
}
