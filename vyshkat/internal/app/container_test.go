package app_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"vyshkat/internal/app"
	"vyshkat/internal/console"
	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
	"vyshkat/internal/inspection"
)

type rejectAll struct{ reason error }

func (r rejectAll) Inspect(domain.Vehicle) error { return r.reason }

func testIO() app.IO {
	return app.IO{In: strings.NewReader(""), Out: &bytes.Buffer{}}
}

func TestContainerResolvesGraph(t *testing.T) {
	c, err := app.BuildContainer(app.DefaultConfig(), testIO())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	err = c.Invoke(func(a *console.App, p *fleet.Park, i fleet.Inspector) {
		if a == nil || p == nil || i == nil {
			t.Error("container returned nil dependency")
		}
		if _, ok := i.(*inspection.ServiceCenter); !ok {
			t.Errorf("default inspector must be ServiceCenter, got %T", i)
		}
	})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
}

func TestContainerParkIsSingleton(t *testing.T) {
	c, err := app.BuildContainer(app.DefaultConfig(), testIO())
	if err != nil {
		t.Fatal(err)
	}
	var first, second *fleet.Park
	_ = c.Invoke(func(p *fleet.Park) { first = p })
	_ = c.Invoke(func(p *fleet.Park) { second = p })
	if first != second {
		t.Error("park must be created once per container")
	}
}

func TestContainerInspectorCanBeReplaced(t *testing.T) {
	c, err := app.BuildContainer(app.DefaultConfig(), testIO())
	if err != nil {
		t.Fatal(err)
	}
	reason := errors.New("inspection station is closed")
	if err := c.Decorate(func(fleet.Inspector) fleet.Inspector { return rejectAll{reason} }); err != nil {
		t.Fatalf("decorate: %v", err)
	}

	s, _ := domain.NewSimplicity(8)
	v, _ := domain.NewBicycle("Stels", s)
	err = c.Invoke(func(p *fleet.Park) {
		_, acceptErr := p.AcceptVehicle(v)
		if !errors.Is(acceptErr, fleet.ErrInspectionFailed) || !errors.Is(acceptErr, reason) {
			t.Errorf("want rejection from replaced inspector, got %v", acceptErr)
		}
	})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
}

func TestContainerUsesConfig(t *testing.T) {
	cfg := app.Config{MaxDailyKWh: 0.1, BeginnerMinSimplicity: 10}
	c, err := app.BuildContainer(cfg, testIO())
	if err != nil {
		t.Fatal(err)
	}
	s, _ := domain.NewSimplicity(9)
	scooter, _ := domain.NewScooter("Ninebot", 0.5, s)
	bicycle, _ := domain.NewBicycle("Stels", s)

	err = c.Invoke(func(p *fleet.Park) {
		if _, err := p.AcceptVehicle(scooter); !errors.Is(err, inspection.ErrExcessiveConsumption) {
			t.Errorf("limit from config must apply, got %v", err)
		}
		if _, err := p.AcceptVehicle(bicycle); err != nil {
			t.Fatal(err)
		}
		if n := len(p.BeginnerVehicles()); n != 0 {
			t.Errorf("beginner threshold from config must apply, got %d vehicles", n)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}
