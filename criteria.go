package daggertech

// Criteria is used to safely build your criteria for searches
type Criteria struct {
	Where      interface{}
	Order      interface{}
	Limit      int
	Offset     int
	IncDeleted bool
}
