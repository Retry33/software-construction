package app

import (
	"io"

	"go.uber.org/dig"

	"vyshkat/internal/console"
	"vyshkat/internal/fleet"
	"vyshkat/internal/inspection"
)

type Config struct {
	MaxDailyKWh           float64
	BeginnerMinSimplicity int
}

func DefaultConfig() Config {
	return Config{
		MaxDailyKWh:           2.0,
		BeginnerMinSimplicity: fleet.DefaultBeginnerMinSimplicity,
	}
}

type IO struct {
	In  io.Reader
	Out io.Writer
}

func BuildContainer(cfg Config, stdio IO) (*dig.Container, error) {
	c := dig.New()
	providers := []any{
		func() Config { return cfg },
		func() IO { return stdio },
		func(cfg Config) fleet.Inspector {
			return inspection.NewServiceCenter(cfg.MaxDailyKWh)
		},
		func(cfg Config) fleet.BeginnerPolicy {
			return fleet.MinSimplicityPolicy{Min: cfg.BeginnerMinSimplicity}
		},
		fleet.NewPark,
		func(p *fleet.Park, stdio IO) *console.App {
			return console.New(p, stdio.In, stdio.Out)
		},
	}
	for _, p := range providers {
		if err := c.Provide(p); err != nil {
			return nil, err
		}
	}
	return c, nil
}
