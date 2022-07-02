package daggertests

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markoxley/daggertech"
	"github.com/markoxley/daggertech/clause"
)

func TestConnection(t *testing.T) {
	conf := getConnectionDetails()
	if !daggertech.Configure(conf) {
		t.Error("Unable to connect to database")
	}
}

func TestNewTableCreation(t *testing.T) {
	reset()
	m := &TestModel{}
	daggertech.Save(m)
	if !testTableExists("TestModel") {
		t.Errorf("Error testing for %v", "TestModel")
	}
}

func TestCount(t *testing.T) {
	reset()
	c1 := daggertech.Count("MISSING_TABLE", nil)
	if c1 != 0 {
		t.Errorf("Expected 0, got %d", c1)
	}
	tm := &TestModel{}
	daggertech.Save(tm)
	daggertech.RawExecute("delete from TestModel")
	tm2 := &TestModel{}
	daggertech.Save(tm2)
	i := daggertech.Count("TestModel", nil)
	if i != 1 {
		t.Errorf("Expected 1, got %d", i)
	}
}

func TestGetRecord(t *testing.T) {
	reset()
	tm1 := &TestModel{
		Name: "Test1",
		Age:  42,
	}
	daggertech.Save(tm1)
	cl := clause.Equal("id", *tm1.ID).ToString()
	c := &daggertech.Criteria{
		Where: cl,
	}
	tm2, _ := daggertech.First(&TestModel{}, c)
	tm3, ok := tm2.(*TestModel)
	if !ok {
		t.Error("Unable to convert Modeller to TestModel")
	}
	if *tm3.ID != *tm1.ID {
		t.Errorf("Expected ID %v, got %v", *tm1.ID, *tm3.ID)
	}
	if compareDates(tm3.CreateDate, tm1.CreateDate) {
		t.Errorf("Expected CreateDate %v, got %v", tm1.CreateDate, tm3.CreateDate)
	}
	if compareDates(tm3.LastUpdate, tm1.LastUpdate) {
		t.Errorf("Expected LastUpdate %v, got %v", tm1.LastUpdate, tm3.LastUpdate)
	}
}

func TestUpdateRecord(t *testing.T) {
	reset()
	tm1 := &TestModel{
		Name: "Test1",
		Age:  42,
	}
	daggertech.Save(tm1)

	tm2, ok := daggertech.First(&TestModel{}, nil)

	if !ok {
		t.Error("Failed to retrieve model")
	}

	tm3, _ := tm2.(*TestModel)

	tm3.Age = 18
	tm3.Name = "David"

	daggertech.Save(tm3)

	i := daggertech.Count("TestModel", nil)
	if i != 1 {
		t.Errorf("Expected 1 record, found %d", i)
	}

	tm4, _ := daggertech.First(&TestModel{}, nil)
	tm5 := tm4.(*TestModel)
	if tm5.Age != tm3.Age {
		t.Errorf("Expected Age of %d, got %d", tm3.Age, tm5.Age)
	}
	if tm5.Name != tm3.Name {
		t.Errorf("Expected Name of %s, got %s", tm3.Name, tm5.Name)
	}
}
