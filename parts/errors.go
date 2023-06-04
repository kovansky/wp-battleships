package parts

import "fmt"

type ErrFieldMalformed struct {
	field string
}

func NewErrFieldMalformed(field string) ErrFieldMalformed {
	return ErrFieldMalformed{field}
}

func (e ErrFieldMalformed) Error() string {
	return fmt.Sprintf("field %s is malformed", e.field)
}

type ErrAdjacentFieldMalformed struct {
	field string
}

func NewErrAdjacentFieldMalformed(field string) ErrAdjacentFieldMalformed {
	return ErrAdjacentFieldMalformed{field: field}
}

func (e ErrAdjacentFieldMalformed) Error() string {
	return fmt.Sprintf("field adjacent to %s is malformed", e.field)
}

type ErrFieldOutOfRange struct{}

func NewErrFieldOutOfRange() ErrFieldOutOfRange {
	return ErrFieldOutOfRange{}
}

func (e ErrFieldOutOfRange) Error() string {
	return "field is out of range"
}

type ErrFieldNonadjacent struct {
	field string
}

func NewErrFieldNonadjacent(field string) ErrFieldNonadjacent {
	return ErrFieldNonadjacent{field: field}
}

func (e ErrFieldNonadjacent) Error() string {
	return fmt.Sprintf("field %s is not adjacent to any of the existing parts of the ship", e.field)
}

type ErrShipSize struct {
	size int
}

func NewErrShipSize(size int) ErrShipSize {
	return ErrShipSize{size: size}
}

func (e ErrShipSize) Error() string {
	return fmt.Sprintf("ship size %d is incorrect (allowed sizes: 1, 2, 3, 4)", e.size)
}
