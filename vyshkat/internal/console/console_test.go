package console_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"vyshkat/internal/console"
	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

type fakeInspector struct{ err error }

func (f fakeInspector) Inspect(domain.Vehicle) error { return f.err }

func run(t *testing.T, insp fleet.Inspector, input string) (*fleet.Park, string) {
	t.Helper()
	park := fleet.NewPark(insp, fleet.DefaultBeginnerPolicy())
	var out bytes.Buffer
	if err := console.New(park, strings.NewReader(input), &out).Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	return park, out.String()
}

func lines(l ...string) string { return strings.Join(l, "\n") + "\n" }

func assertContains(t *testing.T, out string, parts ...string) {
	t.Helper()
	for _, p := range parts {
		if !strings.Contains(out, p) {
			t.Errorf("output does not contain %q\n--- output ---\n%s", p, out)
		}
	}
}

func TestAddVehicleFlow(t *testing.T) {
	park, out := run(t, fakeInspector{}, lines("1", "1", "Ninebot G30", "0,6", "7", "0"))
	assertContains(t, out, "Accepted: Electric scooter Ninebot G30, inventory number VK-0001", "Bye!")
	if len(park.Inventory()) != 1 {
		t.Errorf("vehicle must be added to the park")
	}
}

func TestAddBicycleSkipsEnergyQuestion(t *testing.T) {
	park, out := run(t, fakeInspector{}, lines("1", "2", "Stels", "9", "4", "0"))
	if strings.Contains(out, "Daily energy consumption") {
		t.Error("bicycle must not be asked about energy")
	}
	assertContains(t, out, "Bicycle Stels", "VK-0001")
	if len(park.BeginnerVehicles()) != 1 {
		t.Error("bicycle with simplicity 9 must be suitable for beginners")
	}
}

func TestInspectionRejectionIsShown(t *testing.T) {
	park, out := run(t, fakeInspector{err: errors.New("brakes are worn out")}, lines("1", "3", "Eltreco", "1.2", "8", "0"))
	assertContains(t, out, "Error: vehicle rejected by inspection: brakes are worn out")
	if len(park.Inventory()) != 0 {
		t.Error("rejected vehicle must not be added")
	}
}

func TestRegisterItemFlow(t *testing.T) {
	park, out := run(t, fakeInspector{}, lines("2", "3", "Pokrovka", "6.5", "2", "1", "M", "3", "5", "0"))
	assertContains(t, out,
		"Registered: Charging cabinet Pokrovka, inventory number VK-0001",
		"Registered: Helmet M, inventory number VK-0002",
		"Units on balance:       2",
		"Daily energy, kWh:      6.50",
		"VK-0002  Helmet M",
	)
	if len(park.Inventory()) != 2 {
		t.Errorf("want 2 items, got %d", len(park.Inventory()))
	}
}

func TestInvalidInputIsReported(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"unknown menu option", lines("9", "0"), "Unknown option"},
		{"type out of range", lines("1", "7", "0"), "Error: invalid choice: 7"},
		{"not a number", lines("1", "1", "Ninebot", "abc", "0"), `Error: invalid number: "abc"`},
		{"simplicity out of range", lines("1", "2", "Stels", "11", "0"), "Error: simplicity must be between 1 and 10"},
		{"negative energy", lines("2", "2", "Dock", "-3", "0"), "Error: daily energy consumption must be a non-negative number"},
		{"empty model", lines("2", "1", "", "0"), "Error: model must not be empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			park, out := run(t, fakeInspector{}, tt.input)
			assertContains(t, out, tt.want)
			if len(park.Inventory()) != 0 {
				t.Error("nothing must be added on invalid input")
			}
		})
	}
}

func TestEmptyReports(t *testing.T) {
	_, out := run(t, fakeInspector{}, lines("3", "4", "5", "0"))
	assertContains(t, out, "Units on balance:       0", "No vehicles suitable for beginners.", "The fleet is empty.")
}

func TestEndOfInputExitsGracefully(t *testing.T) {
	_, out := run(t, fakeInspector{}, lines("1", "1", "Ninebot"))
	if strings.Contains(out, "Error") {
		t.Errorf("closed input must not be reported as an error:\n%s", out)
	}
}

func TestDemoCoversAllScenarios(t *testing.T) {
	park := fleet.NewPark(fakeInspector{}, fleet.DefaultBeginnerPolicy())
	var out bytes.Buffer
	console.New(park, strings.NewReader(""), &out).RunDemo()
	assertContains(t, out.String(),
		"accepted as VK-0001",
		"simplicity 12",
		"creation failed",
		"registered as",
		"Fleet summary:",
		"Vehicles suitable for beginners:",
		"Inventory:",
	)
}

func TestMenuRunsDemo(t *testing.T) {
	_, out := run(t, fakeInspector{}, lines("6", "0"))
	assertContains(t, out, "=== Demo scenario ===", "=== Demo finished ===")
}
