// Package clause This is a very basic, and yet versatile ORM package.
// At present, this is only for postgres
package clause

import (
	"fmt"
	"reflect"
)

// Builder is the main clause builder mechanism used for dagger
type Builder struct {
	conjunction conjunction
	children    []clauseInterface
}

func newBuilder(c conjunction) *Builder {
	return &Builder{
		conjunction: c,
		children:    make([]clauseInterface, 0, 5),
	}
}

// ToString returns the string version of the clause
func (c *Builder) ToString() string {

	result := ""
	for _, child := range c.children {
		if result != "" {
			result += string(child.getConjunction())
		}

		v := child.ToString()
		if _, ok := child.(*Builder); ok {
			v = fmt.Sprintf("(%s)", v)
		}
		if v == "" {
			return ""
		}
		result += v
	}
	return result
}

func (c *Builder) getConjunction() conjunction {
	return c.conjunction
}

// Sub creates a new Builder with the first clause being a sub clause
func Sub(c *Builder) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, c)
	return n
}

// Equal creates a new Builder with the first clause being an equal clause
func Equal(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oEqual, false, v))
	return n
}

// Greater creates a new Builder with the first clause being a greater than clause
func Greater(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oGreater, false, v))
	return n
}

// Less creates a new Builder with the first clause being a less than clause
func Less(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLess, false, v))
	return n
}

// Like creates a new Builder with the first clause being a like clause
func Like(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, false, v))
	return n
}

// StartsWith creates a new Builder with the first clause being a starts with clause
func StartsWith(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, false, v+"%"))
	return n
}

// EndsWith creates a new Builder with the first clause being an ends with clause
func EndsWith(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, false, "%"+v))
	return n
}

// Contains creates a new Builder with the first clause being a contains clause
func Contains(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, false, "%"+v+"%"))
	return n
}

// In creates a new Builder with the first clause being an in clause
func In(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oIn, false, values...))
	return n
}

// Between creates a new Builder with the first clause being a between clause
func Between(f string, v1 interface{}, v2 interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oBetween, false, v1, v2))
	return n
}

// IsNull creates a new Builder with the first clause being an is null clause
func IsNull(f string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oIsNull, false))
	return n
}

// IsNotNull creates a new Builder with the first clause being an is null clause
func IsNotNull(f string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oIsNull, true))
	return n
}

// NotEqual creates a new Builder with the first clause being a not equal clause
func NotEqual(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oEqual, true, v))
	return n
}

// NotGreater creates a new Builder with the first clause being a not greater than clause
func NotGreater(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oGreater, true, v))
	return n
}

// NotLess creates a new Builder with the first clause being a not less than clause
func NotLess(f string, v interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLess, true, v))
	return n
}

// NotLike creates a new Builder with the first clause being a not like clause
func NotLike(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, true, v))
	return n
}

// NotStartsWith creates a new Builder with the first clause being a not starts with clause
func NotStartsWith(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, true, v+"%"))
	return n
}

// NotEndsWith creates a new Builder with the first clause being a not ends with clause
func NotEndsWith(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, true, "%"+v))
	return n
}

// NotContains creates a new Builder with the first clause being a not contains clause
func NotContains(f string, v string) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oLike, true, "%"+v+"%"))
	return n
}

// NotIn creates a new Builder with the first clause being a not in clause
func NotIn(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oIn, true, values...))
	return n
}

// NotBetween creates a new Builder with the first clause being a not between clause
func NotBetween(f string, v1 interface{}, v2 interface{}) *Builder {
	n := newBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, f, oBetween, true, v1, v2))
	return n
}

// AndSub add and existing subclause to the clause with an AND conjunction
func (c *Builder) AndSub(n *Builder) *Builder {
	n.conjunction = conAnd
	c.children = append(c.children, n)
	return c
}

// AndEqual add an equal clause to the clause with an AND conjunction
func (c *Builder) AndEqual(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oEqual, false, v))
	return c
}

// AndGreater add a greater than clause to the clause with an AND conjunction
func (c *Builder) AndGreater(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oGreater, false, v))
	return c
}

// AndLess add a less than clause to the clause with an AND conjunction
func (c *Builder) AndLess(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLess, false, v))
	return c
}

// AndLike add a like clause to the clause with an AND conjunction
func (c *Builder) AndLike(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, false, v))
	return c
}

// AndStartsWith add a starts with clause to the clause with an AND conjunction
func (c *Builder) AndStartsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, false, v+"%"))
	return c
}

// AndEndsWith add a ends with clause to the clause with an AND conjunction
func (c *Builder) AndEndsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, false, "%"+v))
	return c
}

// AndContains add a contains clause to the clause with an AND conjunction
func (c *Builder) AndContains(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, false, "%"+v+"%"))
	return c
}

// AndIn add an in clause to the clause with an AND conjunction
func (c *Builder) AndIn(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	c.children = append(c.children, newClause(conAnd, f, oIn, false, values...))
	return c
}

// AndBetween add a between clause to the clause with an AND conjunction
func (c *Builder) AndBetween(f string, v1 interface{}, v2 interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oBetween, false, v1, v2))
	return c
}

// AndIsNull adds an is null clause to the clause with an AND conjunction
func (c *Builder) AndIsNull(f string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oIsNull, false))
	return c
}

// AndNotIsNull adds a not is null clause to the clause with an AND conjunction
func (c *Builder) AndNotIsNull(f string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oIsNull, true))
	return c
}

// AndNotEqual add a not equal clause to the clause with an AND conjunction
func (c *Builder) AndNotEqual(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oEqual, true, v))
	return c
}

// AndNotGreater add a not greater than clause to the clause with an AND conjunction
func (c *Builder) AndNotGreater(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oGreater, true, v))
	return c
}

// AndNotLess add a not less than clause to the clause with an AND conjunction
func (c *Builder) AndNotLess(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLess, true, v))
	return c
}

// AndNotLike add a not like clause to the clause with an AND conjunction
func (c *Builder) AndNotLike(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, true, v))
	return c
}

// AndNotStartsWith add a starts with clause to the clause with an AND conjunction
func (c *Builder) AndNotStartsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, true, v+"%"))
	return c
}

// AndNotEndsWith add a ends with clause to the clause with an AND conjunction
func (c *Builder) AndNotEndsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, true, "%"+v))
	return c
}

// AndNotContains add a contains clause to the clause with an AND conjunction
func (c *Builder) AndNotContains(f string, v string) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oLike, true, "%"+v+"%"))
	return c
}

// AndNotIn add a not in clause to the clause with an AND conjunction
func (c *Builder) AndNotIn(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	c.children = append(c.children, newClause(conAnd, f, oIn, true, values...))
	return c
}

// AndNotBetween add a not between clause to the clause with an AND conjunction
func (c *Builder) AndNotBetween(f string, v1 interface{}, v2 interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, f, oBetween, true, v1, v2))
	return c
}

// OrSub add and existing subclause to the clause with an OR conjunction
func (c *Builder) OrSub(n *Builder) *Builder {
	n.conjunction = conOr
	c.children = append(c.children, n)
	return c
}

// OrEqual add an equal clause to the clause with an OR conjunction
func (c *Builder) OrEqual(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oEqual, false, v))
	return c
}

// OrGreater add a greater than clause to the clause with an OR conjunction
func (c *Builder) OrGreater(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oGreater, false, v))
	return c
}

// OrLess add a less than clause to the clause with an OR conjunction
func (c *Builder) OrLess(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLess, false, v))
	return c
}

// OrLike add a like clause to the clause with an OR conjunction
func (c *Builder) OrLike(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, false, v))
	return c
}

// OrStartsWith add a starts with clause to the clause with an OR conjunction
func (c *Builder) OrStartsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, false, v+"%"))
	return c
}

// OrEndsWith add a ends with clause to the clause with an OR conjunction
func (c *Builder) OrEndsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, false, "%"+v))
	return c
}

// OrContains add a contains clause to the clause with an OR conjunction
func (c *Builder) OrContains(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, false, "%"+v+"%"))
	return c
}

// OrIn add an in clause to the clause with an OR conjunction
func (c *Builder) OrIn(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	c.children = append(c.children, newClause(conOr, f, oIn, false, values...))
	return c
}

// OrBetween add a between clause to the clause with an OR conjunction
func (c *Builder) OrBetween(f string, v1 interface{}, v2 interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oBetween, false, v1, v2))
	return c
}

// OrIsNull adds an is null clause to the clause with an OR conjunction
func (c *Builder) OrIsNull(f string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oIsNull, false))
	return c
}

// OrNotIsNull adds a not is null clause to the clause with an OR conjunction
func (c *Builder) OrNotIsNull(f string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oIsNull, true))
	return c
}

// OrNotEqual add a not equal clause to the clause with an OR conjunction
func (c *Builder) OrNotEqual(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oEqual, true, v))
	return c
}

// OrNotGreater add a not greater than clause to the clause with an OR conjunction
func (c *Builder) OrNotGreater(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oGreater, true, v))
	return c
}

// OrNotLess add a not less than clause to the clause with an OR conjunction
func (c *Builder) OrNotLess(f string, v interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLess, true, v))
	return c
}

// OrNotLike add a not like clause to the clause with an OR conjunction
func (c *Builder) OrNotLike(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, true, v))
	return c
}

// OrNotStartsWith add a not starts with clause to the clause with an OR conjunction
func (c *Builder) OrNotStartsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, true, v+"%"))
	return c
}

// OrNotEndsWith add a not ends with clause to the clause with an OR conjunction
func (c *Builder) OrNotEndsWith(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, true, "%"+v))
	return c
}

// OrNotContains add a not contains clause to the clause with an OR conjunction
func (c *Builder) OrNotContains(f string, v string) *Builder {
	c.children = append(c.children, newClause(conOr, f, oLike, true, "%"+string(v)+"%"))
	return c
}

// OrNotIn add a not in clause to the clause with an OR conjunction
func (c *Builder) OrNotIn(f string, v ...interface{}) *Builder {
	values := consolidateArray(v)
	c.children = append(c.children, newClause(conOr, f, oIn, true, values...))
	return c
}

// OrNotBetween add a not between clause to the clause with an OR conjunction
func (c *Builder) OrNotBetween(f string, v1 interface{}, v2 interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, f, oBetween, true, v1, v2))
	return c
}

func consolidateArray(values []interface{}) []interface{} {
	res := make([]interface{}, 0, len(values))
	for _, v := range values {
		arr := reflect.ValueOf(v)
		if arr.Kind() != reflect.Array {
			res = append(res, v)
			continue
		}
		for i := 0; i < arr.Len(); i++ {
			res = append(res, arr.Index(i).Interface())
		}
	}
	return res
}
