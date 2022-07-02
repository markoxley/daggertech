package daggertech

// Criteria is used to safely build your criteria for searches
type Criteria struct {
	Where      string
	Order      string
	Limit      int
	Offset     int
	IncDeleted bool
}
