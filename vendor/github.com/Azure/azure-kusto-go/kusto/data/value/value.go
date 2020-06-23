/*
Package value holds Kusto data value representations. All types provide a Kusto that
stores the native value and Valid which indicates if the value was set or was null.

Kusto Value

A value.Kusto can hold types that represent Kusto Scalar types that define column data.
We represent that with an interface:

	type Kusto interface

This interface can hold the following values:

	value.Bool
	value.Int
	value.Long
	value.Real
	value.Decimal
	value.String
	value.Dynamic
	value.DateTime
	value.Timespan

Each type defined above has at minimum two fields:

	.Value - The type specific value
	.Valid - True if the value was non-null in the Kusto table

Each provides at minimum the following two methods:

	.String() - Returns the string representation of the value.
	.Unmarshal() - Unmarshals the value into a standard Go type.

The Unmarshal() is for internal use, it should not be needed by an end user. Use .Value or table.Row.ToStruct() instead.
*/
package value

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Kusto represents a Kusto value.
type Kusto interface {
	isKustoVal()
	// String implements fmt.Stringer().
	String() string
}

// Values is a list of Kusto values, usually an ordered row.
type Values []Kusto

// Bool represents a Kusto boolean type. Bool implements Kusto.
type Bool struct {
	// Value holds the value of the type.
	Value bool
	// Valid indicates if this value was set.
	Valid bool
}

func (Bool) isKustoVal() {}

// String implements fmt.Stringer.
func (bo Bool) String() string {
	if !bo.Valid {
		return ""
	}
	if bo.Value {
		return "true"
	}
	return "false"
}

// Unmarshal unmarshals i into Bool. i must be a bool or nil.
func (bo *Bool) Unmarshal(i interface{}) error {
	if i == nil {
		bo.Value = false
		bo.Valid = false
		return nil
	}
	v, ok := i.(bool)
	if !ok {
		return fmt.Errorf("Column with type 'bool' had value that was %T", i)
	}
	bo.Value = v
	bo.Valid = true
	return nil
}

// DateTime represents a Kusto datetime type.  DateTime implements Kusto.
type DateTime struct {
	// Value holds the value of the type.
	Value time.Time
	// Valid indicates if this value was set.
	Valid bool
}

// String implements fmt.Stringer.
func (d DateTime) String() string {
	if !d.Valid {
		return ""
	}
	return fmt.Sprint(d.Value.Format(time.RFC3339Nano))
}

func (DateTime) isKustoVal() {}

// Marshal marshals the DateTime into a Kusto compatible string.
func (d DateTime) Marshal() string {
	if !d.Valid {
		return time.Time{}.Format(time.RFC3339Nano)
	}
	return d.Value.Format(time.RFC3339Nano)
}

// Unmarshal unmarshals i into DateTime. i must be a string representing RFC3339Nano or nil.
func (d *DateTime) Unmarshal(i interface{}) error {
	if i == nil {
		d.Value = time.Time{}
		d.Valid = false
		return nil
	}

	str, ok := i.(string)
	if !ok {
		return fmt.Errorf("Column with type 'datetime' had value that was %T", i)
	}

	t, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		return fmt.Errorf("Column with type 'datetime' had value %s which did not parse: %s", str, err)
	}
	d.Value = t
	d.Valid = true

	return nil
}

// Dynamic represents a Kusto dynamic type.  Dynamic implements Kusto.
type Dynamic struct {
	// Value holds the value of the type.
	Value []byte
	// Valid indicates if this value was set.
	Valid bool
}

func (Dynamic) isKustoVal() {}

// String implements fmt.Stringer.
func (d Dynamic) String() string {
	if !d.Valid {
		return ""
	}

	return string(d.Value)
}

// Unmarshal unmarshal's i into Dynamic. i must be a string, []byte, map[string]interface{}, []interface{}, other JSON serializable value or nil.
// If []byte or string, must be a JSON representation of a value.
func (d *Dynamic) Unmarshal(i interface{}) error {
	if i == nil {
		d.Value = nil
		d.Valid = false
		return nil
	}

	switch v := i.(type) {
	case []byte:
		d.Value = v
		d.Valid = true
		return nil
	case string:
		d.Value = []byte(v)
		d.Valid = true
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("Column with type 'dynamic' was a %T that could not be JSON encoded: %s", i, err)
	}

	d.Value = b
	d.Valid = true
	return nil
}

// GUID represents a Kusto GUID type.  GUID implements Kusto.
type GUID struct {
	// Value holds the value of the type.
	Value uuid.UUID
	// Valid indicates if this value was set.
	Valid bool
}

func (GUID) isKustoVal() {}

// String implements fmt.Stringer.
func (g GUID) String() string {
	if !g.Valid {
		return ""
	}
	return g.Value.String()
}

// Unmarshal unmarshals i into GUID. i must be a string representing a GUID or nil.
func (g *GUID) Unmarshal(i interface{}) error {
	if i == nil {
		g.Value = uuid.UUID{}
		g.Valid = false
		return nil
	}
	str, ok := i.(string)
	if !ok {
		return fmt.Errorf("Column with type 'guid' was not stored as a string, was %T", i)
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("Column with type 'guid' did not store a valid uuid(%s): %s", str, err)
	}
	g.Value = u
	g.Valid = true
	return nil
}

// Int represents a Kusto int type. Values int type's are int32 values.  Int implements Kusto.
type Int struct {
	// Value holds the value of the type.
	Value int32
	// Valid indicates if this value was set.
	Valid bool
}

func (Int) isKustoVal() {}

// String implements fmt.Stringer.
func (in Int) String() string {
	if !in.Valid {
		return ""
	}
	return strconv.Itoa(int(in.Value))
}

// Unmarshal unmarshals i into Int. i must be an int32 or nil.
func (in *Int) Unmarshal(i interface{}) error {
	if i == nil {
		in.Value = 0
		in.Valid = false
		return nil
	}

	var myInt int64

	switch v := i.(type) {
	case json.Number:
		var err error
		myInt, err = v.Int64()
		if err != nil {
			return fmt.Errorf("Column with type 'int' had value json.Number that had error on .Int64(): %s", err)
		}
	case float64:
		if v != math.Trunc(v) {
			return fmt.Errorf("Column with type 'int' had value float64(%v) that did not represent a whole number", v)
		}
		myInt = int64(v)
	case int:
		myInt = int64(v)
	default:
		return fmt.Errorf("Column with type 'int' had value that was not a json.Number or int, was %T", i)
	}

	if myInt > math.MaxInt32 {
		return fmt.Errorf("Column with type 'int' had value that was greater than an int32 can hold, was %d", myInt)
	}
	in.Value = int32(myInt)
	in.Valid = true
	return nil
}

// Long represents a Kusto long type, which is an int64.  Long implements Kusto.
type Long struct {
	// Value holds the value of the type.
	Value int64
	// Valid indicates if this value was set.
	Valid bool
}

func (Long) isKustoVal() {}

// String implements fmt.Stringer.
func (l Long) String() string {
	if !l.Valid {
		return ""
	}
	return strconv.Itoa(int(l.Value))
}

// Unmarshal unmarshals i into Long. i must be an int64 or nil.
func (l *Long) Unmarshal(i interface{}) error {
	if i == nil {
		l.Value = 0
		l.Valid = false
		return nil
	}

	var myInt int64

	switch v := i.(type) {
	case json.Number:
		var err error
		myInt, err = v.Int64()
		if err != nil {
			return fmt.Errorf("Column with type 'long' had value json.Number that had error on .Int64(): %s", err)
		}
	case int:
		myInt = int64(v)
	case float64:
		if v != math.Trunc(v) {
			return fmt.Errorf("Column with type 'int' had value float64(%v) that did not represent a whole number", v)
		}
		myInt = int64(v)
	default:
		return fmt.Errorf("Column with type 'ong' had value that was not a json.Number or int, was %T", i)
	}

	l.Value = myInt
	l.Valid = true
	return nil
}

// Real represents a Kusto real type.  Real implements Kusto.
type Real struct {
	// Value holds the value of the type.
	Value float64
	// Valid indicates if this value was set.
	Valid bool
}

func (Real) isKustoVal() {}

// String implements fmt.Stringer.
func (r Real) String() string {
	if !r.Valid {
		return ""
	}
	return strconv.FormatFloat(r.Value, 'e', -1, 64)
}

// Unmarshal unmarshals i into Real. i must be a json.Number(that is a float64), float64 or nil.
func (r *Real) Unmarshal(i interface{}) error {
	if i == nil {
		r.Value = 0.0
		r.Valid = false
		return nil
	}

	var myFloat float64

	switch v := i.(type) {
	case json.Number:
		var err error
		myFloat, err = v.Float64()
		if err != nil {
			return fmt.Errorf("Column with type 'real' had value json.Number that had error on .Float64(): %s", err)
		}
	case float64:
		myFloat = v
	default:
		return fmt.Errorf("Column with type 'real' had value that was not a json.Number or float64, was %T", i)
	}

	r.Value = myFloat
	r.Valid = true
	return nil
}

// String represents a Kusto string type.  String implements Kusto.
type String struct {
	// Value holds the value of the type.
	Value string
	// Valid indicates if this value was set.
	Valid bool
}

func (String) isKustoVal() {}

// String implements fmt.Stringer.
func (s String) String() string {
	if !s.Valid {
		return ""
	}
	return s.Value
}

// Unmarshal unmarshals i into String. i must be a string or nil.
func (s *String) Unmarshal(i interface{}) error {
	if i == nil {
		s.Value = ""
		s.Valid = false
		return nil
	}

	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("Column with type 'string' had type %T", i)
	}

	s.Value = v
	s.Valid = true
	return nil
}

// Decimal represents a Kusto decimal type.  Decimal implements Kusto.
// Because Go does not have a dynamic decimal type that meets all needs, Decimal
// provides the string representation for you to unmarshal into.
type Decimal struct {
	// Value holds the value of the type.
	Value string
	// Valid indicates if this value was set.
	Valid bool
}

func (Decimal) isKustoVal() {}

// String implements fmt.Stringer.
func (d Decimal) String() string {
	if !d.Valid {
		return ""
	}
	return d.Value
}

// ParseFloat provides builtin support for Go's *big.Float conversion where that type meets your needs.
func (d *Decimal) ParseFloat(base int, prec uint, mode big.RoundingMode) (f *big.Float, b int, err error) {
	return big.ParseFloat(d.Value, base, prec, mode)
}

var decRE = regexp.MustCompile(`^\d*\.\d+$`)

// Unmarshal unmarshals i into Decimal. i must be a string representing a decimal type or nil.
func (d *Decimal) Unmarshal(i interface{}) error {
	if i == nil {
		d.Value = ""
		d.Valid = false
		return nil
	}

	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("Column with type 'decimal' had type %T", i)
	}

	if !decRE.MatchString(v) {
		return fmt.Errorf("column with type 'decimal' does not appear to be a decimal number, was %v", v)
	}

	d.Value = v
	d.Valid = true
	return nil
}

// Timespan represents a Kusto timespan type.  Timespan implements Kusto.
type Timespan struct {
	// Value holds the value of the type.
	Value time.Duration
	// Valid indicates if this value was set.
	Valid bool
}

func (Timespan) isKustoVal() {}

// String implements fmt.Stringer.
func (t Timespan) String() string {
	if !t.Valid {
		return ""
	}
	return t.Value.String()
}

// Marshal marshals the Timespan into a Kusto compatible string.
func (t Timespan) Marshal() string {
	const (
		day = 24 * time.Hour
	)

	if !t.Valid {
		return "00:00:00"
	}

	val := t.Value

	sb := strings.Builder{}
	if t.Value < 0 {
		sb.WriteString("-")
		val = val * -1
	}

	days := val / day
	val = val - (days * day)
	switch {
	case days == 0:
	case days < 10:
		sb.WriteString(fmt.Sprintf("0%d.", int(days)))
	default:
		sb.WriteString(fmt.Sprintf("%d.", int(days)))
	}

	hours := val / time.Hour
	val = val - (hours * time.Hour)
	switch {
	case hours < 10:
		sb.WriteString(fmt.Sprintf("0%d:", int(hours)))
	default:
		sb.WriteString(fmt.Sprintf("%d:", int(hours)))
	}

	minutes := val / time.Minute
	val = val - (minutes * time.Minute)
	switch {
	case minutes < 10:
		sb.WriteString(fmt.Sprintf("0%d:", int(minutes)))
	default:
		sb.WriteString(fmt.Sprintf("%d:", int(minutes)))
	}

	seconds := val / time.Second
	val = val - (seconds * time.Second)
	switch {
	case minutes < 10:
		sb.WriteString(fmt.Sprintf("0%d", int(seconds)))
	default:
		sb.WriteString(fmt.Sprintf("%d", int(seconds)))
	}

	subSecond := strings.Builder{}
	set := false

	milliseconds := val / time.Millisecond
	val = val - (milliseconds * time.Millisecond)
	switch {
	case milliseconds == 0:
		subSecond.WriteString("000")
	case milliseconds < 10:
		set = true
		subSecond.WriteString(fmt.Sprintf("00%d", milliseconds))
	case milliseconds < 100:
		set = true
		subSecond.WriteString(fmt.Sprintf("0%d", milliseconds))
	default:
		set = true
		subSecond.WriteString(fmt.Sprintf("%d", milliseconds))
	}

	nanoseconds := val / time.Nanosecond
	if nanoseconds > 0 {
		set = true
		subSecond.WriteString(fmt.Sprintf("%d", nanoseconds))
	}
	if set {
		sb.WriteString("." + subSecond.String())
	}

	str := strings.TrimRight(sb.String(), "0")
	if strings.HasSuffix(str, ":") {
		str = str + "00"
	}

	return str
}

// Unmarshal unmarshals i into Timespan. i must be a string representing a Values timespan or nil.
func (t *Timespan) Unmarshal(i interface{}) error {
	const (
		hoursIndex   = 0
		minutesIndex = 1
		secondsIndex = 2
	)

	if i == nil {
		t.Value = 0
		t.Valid = false
		return nil
	}

	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("Column with type 'timespan' had type %T", i)
	}

	negative := false
	if len(v) > 1 {
		if string(v[0]) == "-" {
			negative = true
			v = v[1:]
		}
	}

	sp := strings.Split(v, ":")
	if len(sp) != 3 {
		return fmt.Errorf("value to unmarshal into Timespan does not seem to fit format '00:00:00', where values are decimal(%s)", v)
	}

	var sum time.Duration

	d, err := t.unmarshalDaysHours(sp[hoursIndex])
	if err != nil {
		return err
	}
	sum += d

	d, err = t.unmarshalMinutes(sp[minutesIndex])
	if err != nil {
		return err
	}
	sum += d

	d, err = t.unmarshalSeconds(sp[secondsIndex])
	if err != nil {
		return err
	}

	sum += d

	if negative {
		sum = sum * time.Duration(-1)
	}

	t.Value = sum
	t.Valid = true
	return nil
}

var day = 24 * time.Hour

func (t *Timespan) unmarshalDaysHours(s string) (time.Duration, error) {
	sp := strings.Split(s, ".")
	switch len(sp) {
	case 1:
		hours, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("timespan's hours/day field was incorrect, was %s: %s", s, err)
		}
		return time.Duration(hours) * time.Hour, nil
	case 2:
		days, err := strconv.Atoi(sp[0])
		if err != nil {
			return 0, fmt.Errorf("timespan's hours/day field was incorrect, was %s", s)
		}
		hours, err := strconv.Atoi(sp[1])
		if err != nil {
			return 0, fmt.Errorf("timespan's hours/day field was incorrect, was %s", s)
		}
		return time.Duration(days)*day + time.Duration(hours)*time.Hour, nil
	}
	return 0, fmt.Errorf("timespan's hours/days field did not have the requisite '.'s, was %s", s)
}

func (t *Timespan) unmarshalMinutes(s string) (time.Duration, error) {
	s = strings.Split(s, ".")[0] // We can have 01 or 01.00 or 59, but nothing comes behind the .

	minutes, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("timespan's minutes field was incorrect, was %s", s)
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("timespan's minutes field was incorrect, was %s", s)
	}
	return time.Duration(minutes) * time.Minute, nil
}

const tick = 100 * time.Nanosecond

// unmarshalSeconds deals with this crazy output format. Instead of having some multiplier, the number
// of precision characters behind the decimal indicates your multiplier. This can be between 0 and 7, but
// really only has 3, 4 and 7. There is something called a tick, which is 100 Nanoseconds and the precision
// at len 4 is 100 * Microsecond (don't know if that has a name).
func (t *Timespan) unmarshalSeconds(s string) (time.Duration, error) {
	// "03" = 3 * time.Second
	// "00.099" = 99 * time.Millisecond
	// "03.0123" == 3 * time.Second + 12300 * time.Microsecond
	sp := strings.Split(s, ".")
	switch len(sp) {
	case 1:
		seconds, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("timespan's seconds field was incorrect, was %s", s)
		}
		return time.Duration(seconds) * time.Second, nil
	case 2:
		seconds, err := strconv.Atoi(sp[0])
		if err != nil {
			return 0, fmt.Errorf("timespan's seconds field was incorrect, was %s", s)
		}
		n, err := strconv.Atoi(sp[1])
		if err != nil {
			return 0, fmt.Errorf("timespan's seconds field was incorrect, was %s", s)
		}
		var prec time.Duration
		switch len(sp[1]) {
		case 1:
			prec = time.Duration(n) * (100 * time.Millisecond)
		case 2:
			prec = time.Duration(n) * (10 * time.Millisecond)
		case 3:
			prec = time.Duration(n) * time.Millisecond
		case 4:
			prec = time.Duration(n) * 100 * time.Microsecond
		case 5:
			prec = time.Duration(n) * 10 * time.Microsecond
		case 6:
			prec = time.Duration(n) * time.Microsecond
		case 7:
			prec = time.Duration(n) * tick
		case 8:
			prec = time.Duration(n) * (10 * time.Nanosecond)
		case 9:
			prec = time.Duration(n) * time.Nanosecond
		default:
			return 0, fmt.Errorf("timespan's seconds field did not have 1-9 numbers after the decimal, had %v", s)
		}

		return time.Duration(seconds)*time.Second + prec, nil
	}
	return 0, fmt.Errorf("timespan's seconds field did not have the requisite '.'s, was %s", s)
}
