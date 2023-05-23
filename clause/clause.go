package clause

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/markoxley/daggertech/utils"
)

// Clause is the basic clause for the Builder
type clause struct {
	conjunction conjunction
	field       string
	not         bool
	op          int
	values      []interface{}
}

// String converts the clause to a string
func (c *clause) String() string {
	// If opcode is out of range, return error
	if c.op > oIsNull || c.op < oEqual {
		return ""
	}
	opCode := c.op

	// If this is a not clause, update the opcode
	if c.not {
		opCode += len(operators) / 2
	}

	// get the number of values required
	fieldCount := 1
	switch c.op {
	case oBetween:
		fieldCount = 2
	case oIn:
		vt := reflect.TypeOf(c.values[0])
		if vt.Kind() == reflect.Slice || vt.Kind() == reflect.Array {
			ov := c.values[0].([]interface{})
			nv := make([]interface{}, len(ov))
			for i, el := range ov {
				nv[i] = el
			}
			c.values = nv
		}
		fieldCount = len(c.values)
		// edge case, there are no values, return an error
		if fieldCount == 0 {
			return ""
		}
	case oIsNull:
		fieldCount = 0
	}
	// If we do not have enough values, return error
	if len(c.values) < fieldCount {
		return ""
	}

	// need to make sure all values are of same type
	oldType := -1
	vls := make([]string, 0, 5)
	for i := 0; i < fieldCount; i++ {
		if f, t, ok := MakeValue(c.values[i]); ok && (oldType < 0 || oldType == t) {
			vls = append(vls, f)
			oldType = t
		}
	}
	switch c.op {
	case oIn:
		return fmt.Sprintf(operators[opCode], c.field, strings.Join(vls, ","))
	case oBetween:
		v1 := vls[0]
		v2 := vls[1]
		if v1 > v2 {
			v1 = vls[1]
			v2 = vls[0]
		}
		return fmt.Sprintf(operators[opCode], c.field, v1, v2)
	case oIsNull:
		return fmt.Sprintf(operators[opCode], c.field)
	default:
		return fmt.Sprintf(operators[opCode], c.field, vls[0])
	}

}

func (c *clause) getConjunction() conjunction {
	return c.conjunction
}

func newClause(c conjunction, f string, o int, n bool, v ...interface{}) *clause {
	return &clause{
		conjunction: c,
		field:       f,
		not:         n,
		op:          o,
		values:      v,
	}
}

// MakeValue returns a safe string representation of the value, and the type
func MakeValue(v interface{}) (string, int, bool) {
	if f, ok := v.(float32); ok {
		r := fmt.Sprintf("%f", f)
		return r[:len(r)-2], dFloat, true
	}

	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%f", f), dDouble, true
	}

	if f, ok := v.(int); ok {
		return fmt.Sprintf("%d", f), dInt, true
	}
	if f, ok := v.(int32); ok {
		return fmt.Sprintf("%d", f), dInt, true
	}

	if f, ok := v.(int64); ok {
		return fmt.Sprintf("%d", f), dLong, true
	}

	if f, ok := v.(bool); ok {
		if f {
			return "1", dBool, true
		}
		return "0", dBool, true
	}

	if f, ok := v.(string); ok {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(f, "'", "''")), dText, true
	}

	if f, ok := v.(time.Time); ok {
		return fmt.Sprintf("'%s'", utils.TimeToSQL(&f)), dDate, true
	}

	if f, ok := v.(float32); ok {
		r := fmt.Sprintf("%f", f)
		return r[:len(r)-2], dFloat, true
	}

	if f, ok := v.(*float64); ok {
		if f != nil {
			return fmt.Sprintf("%f", *f), dDouble, true
		} else {
			return "NULL", dDouble, true
		}
	}

	if f, ok := v.(*int); ok {
		if f != nil {
			return fmt.Sprintf("%d", *f), dInt, true
		} else {
			return "NULL", dInt, true
		}
	}
	if f, ok := v.(*int32); ok {
		if f != nil {
			return fmt.Sprintf("%d", *f), dInt, true
		} else {
			return "NULL", dInt, true
		}
	}

	if f, ok := v.(*int64); ok {
		if f != nil {
			return fmt.Sprintf("%d", *f), dLong, true
		} else {
			return "NULL", dLong, true
		}
	}

	if f, ok := v.(*bool); ok {
		if f != nil {
			if *f {
				return "1", dBool, true
			}
			return "0", dBool, true
		} else {
			return "NULL", dBool, true
		}
	}

	if f, ok := v.(*string); ok {
		if f != nil {
			return fmt.Sprintf("'%s'", strings.ReplaceAll(*f, "'", "''")), dText, true
		} else {
			return "NULL", dText, true
		}
	}

	if f, ok := v.(*time.Time); ok {
		if f != nil {
			return fmt.Sprintf("'%s'", utils.TimeToSQL(f)), dDate, true
		} else {
			return "NULL", dDate, true
		}
	}

	return "", 0, false
}
