package fleet_test

import (
	"errors"
	"math"
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

type fakeInspector struct {
	err       error
	inspected []domain.Vehicle
}

func (f *fakeInspector) Inspect(v domain.Vehicle) error {
	f.inspected = append(f.inspected, v)
	return f.err
}

type stubPolicy bool

func (p stubPolicy) Allows(domain.Vehicle) bool { return bool(p) }

func simplicity(t *testing.T, v int) domain.Simplicity {
	t.Helper()
	s, err := domain.NewSimplicity(v)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func scooter(t *testing.T, kwh float64, s int) *domain.Scooter {
	t.Helper()
	v, err := domain.NewScooter("Ninebot", kwh, simplicity(t, s))
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func bicycle(t *testing.T, s int) *domain.Bicycle {
	t.Helper()
	v, err := domain.NewBicycle("Stels", simplicity(t, s))
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func newPark(insp fleet.Inspector) *fleet.Park {
	return fleet.NewPark(insp, fleet.DefaultBeginnerPolicy())
}

func TestAcceptVehicleApproved(t *testing.T) {
	insp := &fakeInspector{}
	park := newPark(insp)
	v := scooter(t, 0.5, 7)

	number, err := park.AcceptVehicle(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if number != "VK-0001" || v.Number() != "VK-0001" {
		t.Errorf("want number VK-0001, got %q (item: %q)", number, v.Number())
	}
	if len(insp.inspected) != 1 || insp.inspected[0] != v {
		t.Errorf("inspector must be called exactly once with the vehicle")
	}
	if got := park.Inventory(); len(got) != 1 {
		t.Errorf("want 1 item in inventory, got %d", len(got))
	}
}

func TestAcceptVehicleRejected(t *testing.T) {
	reason := errors.New("brakes are worn out")
	insp := &fakeInspector{err: reason}
	park := newPark(insp)
	v := scooter(t, 0.5, 7)

	_, err := park.AcceptVehicle(v)
	if !errors.Is(err, fleet.ErrInspectionFailed) || !errors.Is(err, reason) {
		t.Fatalf("want ErrInspectionFailed wrapping the reason, got %v", err)
	}
	if v.Number() != "" {
		t.Errorf("rejected vehicle must not get a number, got %q", v.Number())
	}
	if len(park.Inventory()) != 0 {
		t.Errorf("rejected vehicle must not be in inventory")
	}

	insp.err = nil
	number, err := park.AcceptVehicle(scooter(t, 0.5, 7))
	if err != nil || number != "VK-0001" {
		t.Errorf("rejection must not consume a number: got %q, %v", number, err)
	}
}

func TestAcceptVehicleTwice(t *testing.T) {
	park := newPark(&fakeInspector{})
	v := scooter(t, 0.5, 7)
	if _, err := park.AcceptVehicle(v); err != nil {
		t.Fatal(err)
	}
	if _, err := park.AcceptVehicle(v); !errors.Is(err, domain.ErrNumberAlreadyAssigned) {
		t.Errorf("want ErrNumberAlreadyAssigned, got %v", err)
	}
	if len(park.Inventory()) != 1 {
		t.Errorf("vehicle must be stored once")
	}
}

func TestAcceptNilVehicle(t *testing.T) {
	insp := &fakeInspector{}
	park := newPark(insp)
	if _, err := park.AcceptVehicle(nil); !errors.Is(err, fleet.ErrNilItem) {
		t.Errorf("want ErrNilItem, got %v", err)
	}
	if len(insp.inspected) != 0 {
		t.Errorf("inspector must not be called for nil")
	}
}

func TestRegisterItemSkipsInspection(t *testing.T) {
	insp := &fakeInspector{err: errors.New("must not be called")}
	park := newPark(insp)
	helmet, _ := domain.NewHelmet("Helmet M")

	number, err := park.RegisterItem(helmet)
	if err != nil || number != "VK-0001" {
		t.Fatalf("want VK-0001, got %q, %v", number, err)
	}
	if len(insp.inspected) != 0 {
		t.Errorf("items must not go through inspection")
	}
}

func TestRegisterItemRejectsVehicle(t *testing.T) {
	park := newPark(&fakeInspector{})
	if _, err := park.RegisterItem(bicycle(t, 9)); !errors.Is(err, fleet.ErrVehicleRequiresInspection) {
		t.Errorf("want ErrVehicleRequiresInspection, got %v", err)
	}
	if _, err := park.RegisterItem(nil); !errors.Is(err, fleet.ErrNilItem) {
		t.Errorf("want ErrNilItem, got %v", err)
	}
}

func TestNumbersAreSequential(t *testing.T) {
	park := newPark(&fakeInspector{})
	helmet, _ := domain.NewHelmet("Helmet M")
	n1, _ := park.AcceptVehicle(scooter(t, 0.5, 7))
	n2, _ := park.RegisterItem(helmet)
	n3, _ := park.AcceptVehicle(bicycle(t, 8))
	if n1 != "VK-0001" || n2 != "VK-0002" || n3 != "VK-0003" {
		t.Errorf("want VK-0001..VK-0003, got %s %s %s", n1, n2, n3)
	}
}

func TestTotalDailyEnergyAndSummary(t *testing.T) {
	park := newPark(&fakeInspector{})
	helmet, _ := domain.NewHelmet("Helmet M")
	dock, _ := domain.NewDockStation("Dock", 3)
	mustAccept(t, park, scooter(t, 0.5, 7))
	mustAccept(t, park, bicycle(t, 8))
	mustRegister(t, park, helmet)
	mustRegister(t, park, dock)

	if got := park.TotalDailyEnergyKWh(); math.Abs(got-3.5) > 1e-9 {
		t.Errorf("want 3.5 kWh, got %v", got)
	}
	want := fleet.Summary{Total: 4, Vehicles: 2, Items: 2, DailyEnergyKWh: 3.5}
	if got := park.Summary(); got != want {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func TestEmptyPark(t *testing.T) {
	park := newPark(&fakeInspector{})
	if got := park.Summary(); got != (fleet.Summary{}) {
		t.Errorf("empty park summary must be zero, got %+v", got)
	}
	if len(park.BeginnerVehicles()) != 0 || len(park.Inventory()) != 0 {
		t.Errorf("empty park must return empty lists")
	}
}

func TestBeginnerVehiclesDefaultRule(t *testing.T) {
	park := newPark(&fakeInspector{})
	hard := scooter(t, 0.5, 5)
	edge := scooter(t, 0.5, 6)
	easy := bicycle(t, 9)
	helmet, _ := domain.NewHelmet("Helmet M")
	mustAccept(t, park, hard)
	mustAccept(t, park, edge)
	mustAccept(t, park, easy)
	mustRegister(t, park, helmet)

	got := park.BeginnerVehicles()
	if len(got) != 2 || got[0] != edge || got[1] != easy {
		t.Errorf("want vehicles with simplicity >= 6 only, got %v", got)
	}
}

func TestBeginnerVehiclesCustomPolicy(t *testing.T) {
	tests := []struct {
		name   string
		policy fleet.BeginnerPolicy
		want   int
	}{
		{"stricter threshold", fleet.MinSimplicityPolicy{Min: 9}, 1},
		{"stub allows all", stubPolicy(true), 3},
		{"stub allows none", stubPolicy(false), 0},
	}
	for _, tt := range tests {
		park := fleet.NewPark(&fakeInspector{}, tt.policy)
		mustAccept(t, park, scooter(t, 0.5, 5))
		mustAccept(t, park, scooter(t, 0.5, 7))
		mustAccept(t, park, bicycle(t, 9))
		if got := len(park.BeginnerVehicles()); got != tt.want {
			t.Errorf("%s: want %d vehicles, got %d", tt.name, tt.want, got)
		}
	}
}

func TestInventoryReturnsCopy(t *testing.T) {
	park := newPark(&fakeInspector{})
	mustAccept(t, park, scooter(t, 0.5, 7))
	items := park.Inventory()
	items[0] = nil
	if park.Inventory()[0] == nil {
		t.Errorf("Inventory must return a copy")
	}
}

func mustAccept(t *testing.T, p *fleet.Park, v domain.Vehicle) {
	t.Helper()
	if _, err := p.AcceptVehicle(v); err != nil {
		t.Fatal(err)
	}
}

func mustRegister(t *testing.T, p *fleet.Park, item domain.InventoryItem) {
	t.Helper()
	if _, err := p.RegisterItem(item); err != nil {
		t.Fatal(err)
	}
}
