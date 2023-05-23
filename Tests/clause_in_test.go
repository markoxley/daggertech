package daggertests

import (
	"testing"

	"github.com/markoxley/daggertech/clause"
)

func TestClauseIn2(t *testing.T) {
	expected := []string{
		"`ID` in (1,2,3,4)",
		"`Name` in ('Mark','Sally','Oliver')",
		"`Age` in (42)",
		"`Colour` in ('RED')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	sd := 42
	ss := "RED"
	result := clause.In("ID", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = clause.In("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
	result = clause.In("Age", sd).String()
	if result != expected[2] {
		t.Errorf("expecting '%s' got '%s'", expected[2], result)
	}
	result = clause.In("Colour", ss).String()
	if result != expected[3] {
		t.Errorf("expecting '%s' got '%s'", expected[3], result)
	}

}

func TestClauseNotIn2(t *testing.T) {
	expected := []string{
		"`ID` not in (1,2,3,4)",
		"`Name` not in ('Mark','Sally','Oliver')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	result := clause.NotIn("ID", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = clause.NotIn("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}

func TestClauseAndIn2(t *testing.T) {
	expected := []string{
		"`ID` = 2 AND `Size` in (2,4,6)",
		"`ID` = 3 AND `Name` in ('Mark','Sally','Oliver')",
	}
	d := []int{2, 4, 6}
	s := []string{"Mark", "Sally", "Oliver"}
	result := clause.Equal("ID", 2).AndIn("Size", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = clause.Equal("ID", 3).AndIn("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}
