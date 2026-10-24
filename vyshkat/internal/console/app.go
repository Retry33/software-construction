package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

var (
	errInputClosed   = errors.New("input closed")
	errInvalidChoice = errors.New("invalid choice")
	errInvalidNumber = errors.New("invalid number")
	errEmptyModel    = errors.New("model must not be empty")
)

type App struct {
	park *fleet.Park
	in   *bufio.Scanner
	out  io.Writer
}

func New(park *fleet.Park, in io.Reader, out io.Writer) *App {
	return &App{park: park, in: bufio.NewScanner(in), out: out}
}

func (a *App) Run() error {
	fmt.Fprintln(a.out, "VyshKat fleet accounting. Get there before the deadline!")
	for {
		a.printMenu()
		choice, err := a.ask("> ")
		if err != nil {
			return nil
		}
		switch choice {
		case "1":
			err = a.addVehicle()
		case "2":
			err = a.registerItem()
		case "3":
			a.printSummary()
		case "4":
			a.printBeginners()
		case "5":
			a.printInventory()
		case "6":
			a.RunDemo()
		case "0":
			fmt.Fprintln(a.out, "Bye!")
			return nil
		default:
			fmt.Fprintln(a.out, "Unknown option, try again.")
		}
		if errors.Is(err, errInputClosed) {
			return nil
		}
		if err != nil {
			fmt.Fprintf(a.out, "Error: %v\n", err)
		}
	}
}

func (a *App) printMenu() {
	fmt.Fprintln(a.out)
	fmt.Fprintln(a.out, "1. Accept a vehicle (with inspection)")
	fmt.Fprintln(a.out, "2. Register an item")
	fmt.Fprintln(a.out, "3. Fleet summary")
	fmt.Fprintln(a.out, "4. Vehicles for beginners")
	fmt.Fprintln(a.out, "5. Inventory list")
	fmt.Fprintln(a.out, "6. Run demo scenario")
	fmt.Fprintln(a.out, "0. Exit")
}

func (a *App) addVehicle() error {
	labels := make([]string, len(vehicleKinds))
	for i, k := range vehicleKinds {
		labels[i] = k.label
	}
	idx, err := a.choose("Vehicle type", labels)
	if err != nil {
		return err
	}
	kind := vehicleKinds[idx]

	model, err := a.askModel()
	if err != nil {
		return err
	}
	var kwh float64
	if kind.usesEnergy {
		if kwh, err = a.askFloat("Daily energy consumption, kWh: "); err != nil {
			return err
		}
	}
	value, err := a.askInt(fmt.Sprintf("Simplicity for beginners (%d-%d): ", domain.MinSimplicity, domain.MaxSimplicity))
	if err != nil {
		return err
	}
	simplicity, err := domain.NewSimplicity(value)
	if err != nil {
		return err
	}

	name := kind.label + " " + model
	vehicle, err := kind.build(name, kwh, simplicity)
	if err != nil {
		return err
	}
	number, err := a.park.AcceptVehicle(vehicle)
	if err != nil {
		return err
	}
	fmt.Fprintf(a.out, "Accepted: %s, inventory number %s\n", name, number)
	return nil
}

func (a *App) registerItem() error {
	labels := make([]string, len(itemKinds))
	for i, k := range itemKinds {
		labels[i] = k.label
	}
	idx, err := a.choose("Item type", labels)
	if err != nil {
		return err
	}
	kind := itemKinds[idx]

	model, err := a.askModel()
	if err != nil {
		return err
	}
	var kwh float64
	if kind.usesEnergy {
		if kwh, err = a.askFloat("Daily energy consumption, kWh: "); err != nil {
			return err
		}
	}

	name := kind.label + " " + model
	item, err := kind.build(name, kwh)
	if err != nil {
		return err
	}
	number, err := a.park.RegisterItem(item)
	if err != nil {
		return err
	}
	fmt.Fprintf(a.out, "Registered: %s, inventory number %s\n", name, number)
	return nil
}

func (a *App) printSummary() {
	s := a.park.Summary()
	fmt.Fprintln(a.out, "Fleet summary:")
	fmt.Fprintf(a.out, "  Units on balance:       %d\n", s.Total)
	fmt.Fprintf(a.out, "  Vehicles:               %d\n", s.Vehicles)
	fmt.Fprintf(a.out, "  Items:                  %d\n", s.Items)
	fmt.Fprintf(a.out, "  Daily energy, kWh:      %.2f\n", s.DailyEnergyKWh)
}

func (a *App) printBeginners() {
	vehicles := a.park.BeginnerVehicles()
	if len(vehicles) == 0 {
		fmt.Fprintln(a.out, "No vehicles suitable for beginners.")
		return
	}
	fmt.Fprintln(a.out, "Vehicles suitable for beginners:")
	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NUMBER\tNAME\tSIMPLICITY")
	for _, v := range vehicles {
		number, name := "-", "-"
		if item, ok := v.(domain.InventoryItem); ok {
			number, name = item.Number(), item.Name()
		}
		fmt.Fprintf(w, "  %s\t%s\t%d\n", number, name, v.Simplicity().Value())
	}
	w.Flush()
}

func (a *App) printInventory() {
	items := a.park.Inventory()
	if len(items) == 0 {
		fmt.Fprintln(a.out, "The fleet is empty.")
		return
	}
	fmt.Fprintln(a.out, "Inventory:")
	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NUMBER\tNAME\tCATEGORY\tDAILY kWh\tSIMPLICITY")
	for _, item := range items {
		category, simplicity, energy := "item", "-", "-"
		if v, ok := item.(domain.Vehicle); ok {
			category = "vehicle"
			simplicity = strconv.Itoa(v.Simplicity().Value())
		}
		if c, ok := item.(domain.EnergyConsumer); ok {
			energy = fmt.Sprintf("%.2f", c.DailyEnergyKWh())
		}
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n", item.Number(), item.Name(), category, energy, simplicity)
	}
	w.Flush()
}

func (a *App) ask(prompt string) (string, error) {
	fmt.Fprint(a.out, prompt)
	if !a.in.Scan() {
		return "", errInputClosed
	}
	return strings.TrimSpace(a.in.Text()), nil
}

func (a *App) askModel() (string, error) {
	model, err := a.ask("Model: ")
	if err != nil {
		return "", err
	}
	if model == "" {
		return "", errEmptyModel
	}
	return model, nil
}

func (a *App) askInt(prompt string) (int, error) {
	s, err := a.ask(prompt)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errInvalidNumber, s)
	}
	return v, nil
}

func (a *App) askFloat(prompt string) (float64, error) {
	s, err := a.ask(prompt)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errInvalidNumber, s)
	}
	return v, nil
}

func (a *App) choose(title string, options []string) (int, error) {
	fmt.Fprintf(a.out, "%s:\n", title)
	for i, o := range options {
		fmt.Fprintf(a.out, "  %d. %s\n", i+1, o)
	}
	n, err := a.askInt("> ")
	if err != nil {
		return 0, err
	}
	if n < 1 || n > len(options) {
		return 0, fmt.Errorf("%w: %d", errInvalidChoice, n)
	}
	return n - 1, nil
}
