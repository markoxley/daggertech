package clause

type clauseInterface interface {
	ToString() string
	getConjunction() conjunction
}
