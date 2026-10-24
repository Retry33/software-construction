package console

import (
	"fmt"

	"vyshkat/internal/domain"
)

func (a *App) RunDemo() {
	fmt.Fprintln(a.out, "=== Demo scenario ===")

	a.step("Accepting vehicles through inspection")
	a.demoAccept(asVehicle(domain.NewScooter("Electric scooter Ninebot G30", 0.6, a.simplicity(7))))
	a.demoAccept(asVehicle(domain.NewBicycle("Bicycle Stels Navigator", a.simplicity(9))))
	a.demoAccept(asVehicle(domain.NewElectricBike("Electric bike Eltreco Ultra", 1.2, a.simplicity(4))))
	a.demoAccept(asVehicle(domain.NewScooter("Electric scooter Kugoo Kirin", 0.9, a.simplicity(6))))

	a.step("Inspection rejects a vehicle with excessive energy consumption")
	a.demoAccept(asVehicle(domain.NewElectricBike("Electric bike Prototype X", 4.5, a.simplicity(8))))

	a.step("Invalid input is rejected before reaching the fleet")
	if _, err := domain.NewSimplicity(12); err != nil {
		fmt.Fprintf(a.out, "  simplicity 12: %v\n", err)
	}
	a.demoAccept(asVehicle(domain.NewScooter("Electric scooter Broken", -1, a.simplicity(5))))

	a.step("Registering items (no inspection required)")
	a.demoRegister(asItem(domain.NewHelmet("Helmet size M")))
	a.demoRegister(asItem(domain.NewHelmet("Helmet size L")))
	a.demoRegister(asItem(domain.NewDockStation("Dock station Pokrovka", 3.0)))
	a.demoRegister(asItem(domain.NewChargingCabinet("Charging cabinet Myasnitskaya", 6.5)))

	a.step("Reports")
	a.printSummary()
	fmt.Fprintln(a.out)
	a.printBeginners()
	fmt.Fprintln(a.out)
	a.printInventory()
	fmt.Fprintln(a.out, "=== Demo finished ===")
}

func (a *App) step(title string) {
	fmt.Fprintf(a.out, "\n-- %s --\n", title)
}

func (a *App) simplicity(v int) domain.Simplicity {
	s, _ := domain.NewSimplicity(v)
	return s
}

func (a *App) demoAccept(v domain.Vehicle, err error) {
	if err != nil {
		fmt.Fprintf(a.out, "  creation failed: %v\n", err)
		return
	}
	name := v.(domain.InventoryItem).Name()
	number, err := a.park.AcceptVehicle(v)
	if err != nil {
		fmt.Fprintf(a.out, "  %s: %v\n", name, err)
		return
	}
	fmt.Fprintf(a.out, "  %s: accepted as %s\n", name, number)
}

func (a *App) demoRegister(item domain.InventoryItem, err error) {
	if err != nil {
		fmt.Fprintf(a.out, "  creation failed: %v\n", err)
		return
	}
	number, err := a.park.RegisterItem(item)
	if err != nil {
		fmt.Fprintf(a.out, "  %s: %v\n", item.Name(), err)
		return
	}
	fmt.Fprintf(a.out, "  %s: registered as %s\n", item.Name(), number)
}
