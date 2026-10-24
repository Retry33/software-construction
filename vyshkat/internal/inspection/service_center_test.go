package inspection_test

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
	"vyshkat/internal/inspection"
)

var _ fleet.Inspector = (*inspection.ServiceCenter)(nil)

func TestServiceCenterInspect(t *testing.T) {
	s, _ := domain.NewSimplicity(7)
	normal, _ := domain.NewScooter("Ninebot", 0.5, s)
	hungry, _ := domain.NewElectricBike("Prototype", 5, s)
	bicycle, _ := domain.NewBicycle("Stels", s)

	center := inspection.NewServiceCenter(1.0)

	tests := []struct {
		name    string
		vehicle domain.Vehicle
		wantErr error
	}{
		{"consumption within limit", normal, nil},
		{"consumption above limit", hungry, inspection.ErrExcessiveConsumption},
		{"non-consumer passes energy check", bicycle, nil},
		{"nil vehicle", nil, inspection.ErrNilVehicle},
	}
	for _, tt := range tests {
		err := center.Inspect(tt.vehicle)
		if tt.wantErr == nil && err != nil {
			t.Errorf("%s: unexpected error %v", tt.name, err)
		}
		if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
			t.Errorf("%s: want %v, got %v", tt.name, tt.wantErr, err)
		}
	}
}
