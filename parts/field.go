package parts

import (
	"strconv"
	"strings"
)

type Field struct {
	identifier string
	numeric    uint8
	adjacent   map[string]string
}

func NewField(s string) (Field, error) {
	var err error
	field := Field{identifier: s}

	field.numeric, err = field.calculateNumeric()
	if err != nil {
		return Field{}, NewErrFieldMalformed(s)
	}

	field.adjacent, err = field.calculateAdjacent()
	if err != nil {
		return Field{}, NewErrAdjacentFieldMalformed(field.String())
	}

	return field, nil
}

func (f Field) Numeric() uint8 {
	return f.numeric
}

func (f Field) calculateNumeric() (uint8, error) {
	s := strings.ToUpper(f.String())
	letter := s[0] - 'A'
	number, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, err
	}

	return (letter * 10) + uint8(number) - 1, nil
}

func (f Field) String() string {
	return f.identifier
}

func (f Field) Adjacent() map[string]string {
	return f.adjacent
}

func (f Field) calculateAdjacent() (map[string]string, error) {
	var (
		err                    error
		hasN, hasS, hasW, hasE bool
	)
	const (
		NS = 1
		EW = 10
	)

	adjacent := make(map[string]string)

	if f.numeric >= 10 {
		adjacent["W"], err = NumericToIdentifier(f.numeric - EW)
		hasW = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric < 90 {
		adjacent["E"], err = NumericToIdentifier(f.numeric + EW)
		hasE = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric%10 > 0 {
		adjacent["S"], err = NumericToIdentifier(f.numeric - NS)
		hasS = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric%10 < 9 {
		adjacent["N"], err = NumericToIdentifier(f.numeric + NS)
		hasN = true
		if err != nil {
			return nil, err
		}
	}

	if hasN && hasW {
		adjacent["NW"], err = NumericToIdentifier(f.numeric + NS - EW)
		if err != nil {
			return nil, err
		}
	}
	if hasN && hasE {
		adjacent["NE"], err = NumericToIdentifier(f.numeric + NS + EW)
		if err != nil {
			return nil, err
		}
	}
	if hasS && hasW {
		adjacent["SW"], err = NumericToIdentifier(f.numeric - NS - EW)
		if err != nil {
			return nil, err
		}
	}
	if hasS && hasE {
		adjacent["SE"], err = NumericToIdentifier(f.numeric - NS + EW)
		if err != nil {
			return nil, err
		}
	}

	for _, f := range adjacent {
		adjacent[f] = f
	}

	return adjacent, nil
}

func (f Field) IsAdjacent(location string) bool {
	_, ok := f.adjacent[location]
	return ok
}

func (f Field) IsCorner(location string) bool {
	if !f.IsAdjacent(location) {
		return false
	}

	target := f.adjacent[location]
	switch target {
	case f.adjacent["SW"],
		f.adjacent["NW"],
		f.adjacent["SE"],
		f.adjacent["NE"]:
		return true
	default:
		return false
	}
}

func (f Field) State(state State) StatedField {
	return StatedField{Field: f, State: state}
}

func NumericToIdentifier(numeric uint8) (string, error) {
	if numeric > 99 {
		return "", NewErrFieldOutOfRange()
	}

	var identifier strings.Builder
	identifier.WriteRune(rune((numeric/10)%10 + 'A'))
	identifier.WriteString(strconv.Itoa(int(numeric%10) + 1))

	return identifier.String(), nil
}

func IsFieldIdentifier(s string) bool {
	if !(s[0] >= 'A' && s[0] <= 'J') {
		return false
	}

	num, err := strconv.Atoi(s[1:])
	if err != nil {
		return false
	}

	if !(num > 0 && num <= 10) {
		return false
	}

	return true
}

type StatedField struct {
	Field
	State State
}

func NewStatedField(s string, state State) (StatedField, error) {
	f, err := NewField(s)
	if err != nil {
		return StatedField{}, err
	}

	return f.State(state), nil
}
