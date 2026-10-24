package fleet

import (
	"errors"
	"fmt"

	"vyshkat/internal/domain"
)

var (
	ErrNilItem                   = errors.New("item must not be nil")
	ErrInspectionFailed          = errors.New("vehicle rejected by inspection")
	ErrNotInventoryItem          = errors.New("object cannot be put on the balance")
	ErrVehicleRequiresInspection = errors.New("vehicles must be accepted through inspection")
)

type Summary struct {
	Total          int
	Vehicles       int
	Items          int
	DailyEnergyKWh float64
}

type Park struct {
	inspector Inspector
	beginners BeginnerPolicy
	items     []domain.InventoryItem
	issued    int
}

func NewPark(inspector Inspector, beginners BeginnerPolicy) *Park {
	return &Park{inspector: inspector, beginners: beginners}
}

func (p *Park) AcceptVehicle(v domain.Vehicle) (string, error) {
	if v == nil {
		return "", ErrNilItem
	}
	item, ok := v.(domain.InventoryItem)
	if !ok {
		return "", ErrNotInventoryItem
	}
	if err := p.inspector.Inspect(v); err != nil {
		return "", fmt.Errorf("%w: %w", ErrInspectionFailed, err)
	}
	return p.register(item)
}

func (p *Park) RegisterItem(item domain.InventoryItem) (string, error) {
	if item == nil {
		return "", ErrNilItem
	}
	if _, ok := item.(domain.Vehicle); ok {
		return "", ErrVehicleRequiresInspection
	}
	return p.register(item)
}

func (p *Park) register(item domain.InventoryItem) (string, error) {
	assigner, ok := item.(domain.NumberAssigner)
	if !ok {
		return "", ErrNotInventoryItem
	}
	number := fmt.Sprintf("VK-%04d", p.issued+1)
	if err := assigner.AssignNumber(number); err != nil {
		return "", err
	}
	p.issued++
	p.items = append(p.items, item)
	return number, nil
}

func (p *Park) Inventory() []domain.InventoryItem {
	out := make([]domain.InventoryItem, len(p.items))
	copy(out, p.items)
	return out
}

func (p *Park) TotalDailyEnergyKWh() float64 {
	var total float64
	for _, item := range p.items {
		if c, ok := item.(domain.EnergyConsumer); ok {
			total += c.DailyEnergyKWh()
		}
	}
	return total
}

func (p *Park) BeginnerVehicles() []domain.Vehicle {
	var out []domain.Vehicle
	for _, item := range p.items {
		if v, ok := item.(domain.Vehicle); ok && p.beginners.Allows(v) {
			out = append(out, v)
		}
	}
	return out
}

func (p *Park) Summary() Summary {
	s := Summary{Total: len(p.items), DailyEnergyKWh: p.TotalDailyEnergyKWh()}
	for _, item := range p.items {
		if _, ok := item.(domain.Vehicle); ok {
			s.Vehicles++
		} else {
			s.Items++
		}
	}
	return s
}
