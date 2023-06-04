package parts

var AllowedSizes = map[int]interface{}{1: nil, 2: nil, 3: nil, 4: nil}

type Ship struct {
	finished bool
	ship     map[string]Field
}

func NewShip() Ship {
	return Ship{
		finished: false,
		ship:     make(map[string]Field),
	}
}

func (s Ship) Contains(field string) bool {
	_, ok := s.ship[field]
	return ok
}

func (s Ship) Size() int {
	return len(s.ship)
}

func (s Ship) Add(field string) (Ship, error) {
	if _, ok := AllowedSizes[s.Size()+1]; !ok {
		return s, NewErrShipSize(s.Size() + 1)
	}

	f, err := NewField(field)
	if err != nil {
		return s, err
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
			return s, NewErrFieldNonadjacent(field)
		}
	}

	s.ship[f.identifier] = f

	return s, nil
}

func (s Ship) Finish() (Ship, error) {
	if _, ok := AllowedSizes[s.Size()]; !ok {
		return s, NewErrShipSize(s.Size())
	}

	s.finished = true
	return s, nil
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

func (s Ship) IsAdjacent(field string) bool {
	for _, f := range s.ship {
		if f.IsAdjacent(field) {
			return true
		}
	}

	return false
}
