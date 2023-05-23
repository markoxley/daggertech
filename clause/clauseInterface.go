package clause

type clauseInterface interface {
	String() string
	getConjunction() conjunction
}
