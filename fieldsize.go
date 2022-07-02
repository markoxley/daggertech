package daggertech

import "fmt"

type fieldSize struct {
	size    int
	decimal int
}

func newSize(sz, dec int) fieldSize {
	return fieldSize{
		size:    sz,
		decimal: dec,
	}
}

func (s fieldSize) toString() string {
	if s.decimal > 0 {
		return fmt.Sprintf("%d,%d", s.size, s.decimal)
	}
	return fmt.Sprintf("%d", s.size)
}
