package domain_test

import (
	"errors"
	"math"
	"testing"

	"vyshkat/internal/domain"
)

func mustSimplicity(t *testing.T, v int) domain.Simplicity {
	t.Helper()
	s, err := domain.NewSimplicity(v)
	if err != nil {
		t.Fatalf("NewSimplicity(%d): %v", v, err)
	}
	return s
}

func TestNewSimplicity(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{value: 0, wantErr: true},
		{value: -3, wantErr: true},
		{value: 1, wantErr: false},
		{value: 6, wantErr: false},
		{value: 10, wantErr: false},
		{value: 11, wantErr: true},
	}
	for _, tt := range tests {
		s, err := domain.NewSimplicity(tt.value)
		if tt.wantErr {
			if !errors.Is(err, domain.ErrInvalidSimplicity) {
				t.Errorf("NewSimplicity(%d): want ErrInvalidSimplicity, got %v", tt.value, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("NewSimplicity(%d): unexpected error %v", tt.value, err)
		}
		if s.Value() != tt.value {
			t.Errorf("NewSimplicity(%d).Value() = %d", tt.value, s.Value())
		}
	}
}

func TestZeroSimplicityRejected(t *testing.T) {
	_, err := domain.NewBicycle("Stels", domain.Simplicity{})
	if !errors.Is(err, domain.ErrInvalidSimplicity) {
		t.Fatalf("want ErrInvalidSimplicity for zero simplicity, got %v", err)
	}
}

func TestEnergyValidation(t *testing.T) {
	s := mustSimplicity(t, 7)
	for _, kwh := range []float64{-0.1, math.NaN(), math.Inf(1)} {
		if _, err := domain.NewScooter("Ninebot", kwh, s); !errors.Is(err, domain.ErrNegativeEnergy) {
			t.Errorf("NewScooter with %v kWh: want ErrNegativeEnergy, got %v", kwh, err)
		}
	}
	if _, err := domain.NewScooter("Ninebot", 0, s); err != nil {
		t.Errorf("zero consumption must be allowed, got %v", err)
	}
}

func TestEmptyName(t *testing.T) {
	if _, err := domain.NewHelmet("   "); !errors.Is(err, domain.ErrEmptyName) {
		t.Fatalf("want ErrEmptyName, got %v", err)
	}
}

func TestCapabilities(t *testing.T) {
	s := mustSimplicity(t, 6)
	scooter, _ := domain.NewScooter("Ninebot", 0.5, s)
	bicycle, _ := domain.NewBicycle("Stels", s)
	ebike, _ := domain.NewElectricBike("Eltreco", 0.8, s)
	helmet, _ := domain.NewHelmet("Helmet M")
	dock, _ := domain.NewDockStation("Main building dock", 3)
	cabinet, _ := domain.NewChargingCabinet("Cabinet 1", 5)

	tests := []struct {
		name         string
		obj          any
		wantVehicle  bool
		wantConsumer bool
	}{
		{"scooter", scooter, true, true},
		{"bicycle", bicycle, true, false},
		{"electric bike", ebike, true, true},
		{"helmet", helmet, false, false},
		{"dock station", dock, false, true},
		{"charging cabinet", cabinet, false, true},
	}
	for _, tt := range tests {
		if _, ok := tt.obj.(domain.InventoryItem); !ok {
			t.Errorf("%s: must be InventoryItem", tt.name)
		}
		if _, ok := tt.obj.(domain.Vehicle); ok != tt.wantVehicle {
			t.Errorf("%s: Vehicle = %v, want %v", tt.name, ok, tt.wantVehicle)
		}
		if _, ok := tt.obj.(domain.EnergyConsumer); ok != tt.wantConsumer {
			t.Errorf("%s: EnergyConsumer = %v, want %v", tt.name, ok, tt.wantConsumer)
		}
	}
}

func TestAssignNumber(t *testing.T) {
	helmet, err := domain.NewHelmet("Helmet L")
	if err != nil {
		t.Fatal(err)
	}
	if helmet.Number() != "" {
		t.Fatalf("new item must have no number, got %q", helmet.Number())
	}
	if err := helmet.AssignNumber(""); !errors.Is(err, domain.ErrEmptyNumber) {
		t.Errorf("empty number: want ErrEmptyNumber, got %v", err)
	}
	if err := helmet.AssignNumber("VK-0001"); err != nil {
		t.Fatalf("first assignment: %v", err)
	}
	if err := helmet.AssignNumber("VK-0002"); !errors.Is(err, domain.ErrNumberAlreadyAssigned) {
		t.Errorf("second assignment: want ErrNumberAlreadyAssigned, got %v", err)
	}
	if helmet.Number() != "VK-0001" {
		t.Errorf("number changed after rejected assignment: %q", helmet.Number())
	}
}
