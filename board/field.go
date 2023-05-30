package board

import (
	"errors"
	"strconv"
	"strings"
)

type Field struct {
	Identifier string
	numeric    uint8
	adjacent   map[string]string
}

func NewField(s string) (Field, error) {
	var err error
	field := Field{Identifier: s}

	field.numeric, err = field.calculateNumeric()
	if err != nil {
		return Field{}, errors.New("malformed field input")
	}

	field.adjacent, err = field.calculateAdjacent()
	if err != nil {
		return Field{}, errors.New("adjacent calculation: malformed field input")
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
	return f.Identifier
}

func (f Field) Adjacent() map[string]string {
	return f.adjacent
}

func (f Field) calculateAdjacent() (map[string]string, error) {
	var (
		err                error
		isN, isS, isW, isE bool
	)
	adjacent := make(map[string]string)

	if f.numeric > 10 {
		adjacent["N"], err = NumericToIdentifier(f.numeric - 10)
		isN = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric < 90 {
		adjacent["S"], err = NumericToIdentifier(f.numeric + 10)
		isS = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric%10 > 0 {
		adjacent["W"], err = NumericToIdentifier(f.numeric - 1)
		isW = true
		if err != nil {
			return nil, err
		}
	}
	if f.numeric%10 < 9 {
		adjacent["E"], err = NumericToIdentifier(f.numeric + 1)
		isE = true
		if err != nil {
			return nil, err
		}
	}

	if isN && isW {
		adjacent["NW"], err = NumericToIdentifier(f.numeric - 11)
		if err != nil {
			return nil, err
		}
	}
	if isN && isE {
		adjacent["NE"], err = NumericToIdentifier(f.numeric - 9)
		if err != nil {
			return nil, err
		}
	}
	if isS && isW {
		adjacent["SW"], err = NumericToIdentifier(f.numeric + 9)
		if err != nil {
			return nil, err
		}
	}
	if isS && isE {
		adjacent["SE"], err = NumericToIdentifier(f.numeric + 11)
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

func NumericToIdentifier(numeric uint8) (string, error) {
	if numeric > 99 {
		return "", errors.New("field out of range")
	}

	var identifier strings.Builder
	identifier.WriteRune(rune((numeric/10)%10 + 'A'))
	identifier.WriteString(strconv.Itoa(int(numeric%10) + 1))
	
	return identifier.String(), nil
}
