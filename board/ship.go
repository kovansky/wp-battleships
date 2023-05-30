package board

import (
	"errors"
)

var AllowedSizes = map[int]interface{}{1: nil, 2: nil, 3: nil, 4: nil}

type Ship struct {
	finished bool
	ship     map[string]Field
}

func NewShip() Ship {
	return Ship{finished: false}
}

func (s Ship) Contains(field string) bool {
	_, ok := s.ship[field]
	return ok
}

func (s Ship) Size() int {
	return len(s.ship)
}

func (s Ship) Add(field string) error {
	if _, ok := AllowedSizes[s.Size()+1]; !ok {
		return errors.New("ship will grow too big if you do this")
	}

	f, err := NewField(field)
	if err != nil {
		return err
	}

	if s.Size() > 0 {
		anyAdjacent := false
		edges := []string{
			f.adjacent["N"],
			f.adjacent["S"],
			f.adjacent["W"],
			f.adjacent["E"],
		}

		for _, edge := range edges {
			if s.Contains(edge) {
				anyAdjacent = true
				break
			}
		}

		if !anyAdjacent {
			return errors.New("field is not adjacent to any of the existing parts of the ship")
		}
	}

	s.ship[f.Identifier] = f

	return nil
}

func (s Ship) Remove(field string) error {
	if !s.Contains(field) {
		return nil
	}
	f := s.ship[field]

	connections := 0

	edges := []string{
		f.adjacent["N"],
		f.adjacent["S"],
		f.adjacent["W"],
		f.adjacent["E"],
	}

	for _, edge := range edges {
		if s.Contains(edge) {
			connections++
			break
		}
	}

	if connections > 0 {
		return errors.New("removing this field would destroy the whole ship")
	}

	delete(s.ship, field)
	return nil
}

func (s Ship) Finish() error {
	if _, ok := AllowedSizes[s.Size()]; !ok {
		return errors.New("incorrect ship size")
	}

	s.finished = true
	return nil
}

func (s Ship) Ship() map[string]Field {
	return s.ship
}

func (s Ship) Protected() (map[string]StatedField, error) {
	protected := make(map[string]StatedField)

	for _, f := range s.ship {
		for direction, adjacent := range f.Adjacent() {
			if IsFieldIdentifier(direction) {
				continue
			}

			if _, inProtected := protected[adjacent]; s.Contains(adjacent) || inProtected {
				continue
			}

			var err error
			var state State = FieldProtected
			if len(direction) == 2 {
				state = FieldCorner
			}

			protected[adjacent], err = NewStatedField(adjacent, state)
			if err != nil {
				return nil, err
			}
		}
	}

	return protected, nil
}
