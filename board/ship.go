package board

type Ship struct {
	finished bool
	ship map[string]interface{}
}

func CreateShip() Ship {
	return Ship{finished: false}
}

func (s Ship) Contains(field string) bool {
	_, ok := s.ship[field]
	return ok
}
