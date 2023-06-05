package parts

type State string

const (
	FieldEmpty     State = "empty"
	FieldShip            = "ship"
	FieldHit             = "hit"
	FieldMiss            = "miss"
	FieldProtected       = "protected"
	FieldCorner          = "corner"
)

func (s State) Priority() int {
	switch s {
	case FieldEmpty:
		return 0
	case FieldShip:
		return 2
	case FieldHit:
		return 3
	case FieldProtected, FieldCorner:
		return 4
	default:
		return 1
	}
}
