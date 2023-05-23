// Package order This is a very basic, and yet versatile ORM package.
// At present, this is only for postgres
package order

import "fmt"

// Builder is the main order builder mechanism used for dagger
type Builder struct {
	fields []order
}

type order struct {
	field     string
	ascending bool
}

func newOrder(f string, a bool) order {
	return order{
		field:     f,
		ascending: a,
	}
}

func newBuilder() *Builder {
	return &Builder{
		fields: make([]order, 0, 0),
	}
}
func (o *order) String() string {
	d := "asc"
	if !o.ascending {
		d = "desc"
	}
	return fmt.Sprintf("`%s` %s", o.field, d)
}

// Desc creates a new Builder with the first ordering field being descending
func Desc(f string) *Builder {
	b := newBuilder()
	b.fields = append(b.fields, newOrder(f, false))
	return b
}

// Asc creates a new Builder with the first ordering field being ascending
func Asc(f string) *Builder {
	b := newBuilder()
	b.fields = append(b.fields, newOrder(f, true))
	return b
}

// Desc add a field being descending
func (b *Builder) Desc(f string) *Builder {
	b.fields = append(b.fields, newOrder(f, false))
	return b
}

// Asc add a field being ascending
func (b *Builder) Asc(f string) *Builder {
	b.fields = append(b.fields, newOrder(f, true))
	return b
}

// String returns the string version of the ordering list
func (b *Builder) String() string {
	r := ""
	for _, o := range b.fields {
		if r != "" {
			r += ", "
		}
		r += o.String()
	}
	return r
}
