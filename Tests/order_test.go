package daggertests

import (
	"testing"

	"github.com/markoxley/daggertech/order"
)

func TestOrderDesc(t *testing.T) {
	e := "`Age` desc"
	b := order.Desc("Age").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}

func TestOrderAsc(t *testing.T) {
	e := "`dob` asc"
	b := order.Asc("dob").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}

func TestOrderAscAsc(t *testing.T) {
	e := "`dob` asc, `age` asc"
	b := order.Asc("dob").Asc("age").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}

func TestOrderAscDesc(t *testing.T) {
	e := "`dob` asc, `age` desc"
	b := order.Asc("dob").Desc("age").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}

func TestOrderDescAsc(t *testing.T) {
	e := "`dob` desc, `age` asc"
	b := order.Desc("dob").Asc("age").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}

func TestOrderDescDesc(t *testing.T) {
	e := "`dob` desc, `age` desc"
	b := order.Desc("dob").Desc("age").ToString()

	if e != b {
		t.Errorf("Expected %v, got %v", e, b)
	}
}
