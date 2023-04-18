package ships

import battleships "github.com/kovansky/wp-battleships"

var _ battleships.Player = (*Player)(nil)

type Player struct {
	name        string
	description string
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Description() string {
	return p.description
}
