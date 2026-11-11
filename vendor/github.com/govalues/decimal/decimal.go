package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

// Decimal represents a finite floating-point decimal number.
// Its zero value corresponds to the numeric value of 0.
// Decimal is designed to be safe for concurrent use by multiple goroutines.
type Decimal struct {
	neg   bool // indicates whether the decimal is negative
	scale int8 // position of the floating decimal point
	coef  fint // numeric value without decimal point
}

const (
	MaxPrec  = 19      // MaxPrec is the maximum length of the coefficient in decimal digits.
	MinScale = 0       // MinScale is the minimum number of digits after the decimal point.
	MaxScale = 19      // MaxScale is the maximum number of digits after the decimal point.
	maxCoef  = maxFint // maxCoef is the maximum absolute value of the coefficient, which is equal to (10^MaxPrec - 1).
)

var (
	NegOne              = MustNew(-1, 0)                         // NegOne represents the decimal value of -1.
	Zero                = MustNew(0, 0)                          // Zero represents the decimal value of 0. For comparison purposes, use the IsZero method.
	One                 = MustNew(1, 0)                          // One represents the decimal value of 1.
	Two                 = MustNew(2, 0)                          // Two represents the decimal value of 2.
	Ten                 = MustNew(10, 0)                         // Ten represents the decimal value of 10.
	Hundred             = MustNew(100, 0)                        // Hundred represents the decimal value of 100.
	Thousand            = MustNew(1_000, 0)                      // Thousand represents the decimal value of 1,000.
	E                   = MustNew(2_718_281_828_459_045_235, 18) // E represents Euler’s number rounded to 18 digits.
	Pi                  = MustNew(3_141_592_653_589_793_238, 18) // Pi represents the value of π rounded to 18 digits.
	errDecimalOverflow  = errors.New("decimal overflow")
	errInvalidDecimal   = errors.New("invalid decimal")
	errScaleRange       = errors.New("scale out of range")
	errInvalidOperation = errors.New("invalid operation")
	errInexactDivision  = errors.New("inexact division")
	errDivisionByZero   = errors.New("division by zero")
)

// newUnsafe creates a new decimal without checking the scale and coefficient.
// Use it only if you are absolutely sure that the arguments are valid.
func newUnsafe(neg bool, coef fint, scale int) Decimal {
	if coef == 0 {
		neg = false
	}
	//nolint:gosec
	return Decimal{neg: neg, coef: coef, scale: int8(scale)}
}

// newSafe creates a new decimal and checks the scale and coefficient.
func newSafe(neg bool, coef fint, scale int) (Decimal, error) {
	switch {
	case scale < MinScale || scale > MaxScale:
		return Decimal{}, errScaleRange
	case coef > maxCoef:
		return Decimal{}, errDecimalOverflow
	}
	return newUnsafe(neg, coef, scale), nil
}

// newFromFint creates a new decimal from a uint64 coefficient.
// This method does not use overflowError to return descriptive errors,
// as it must be as fast as possible.
func newFromFint(neg bool, coef fint, scale, minScale int) (Decimal, error) {
	var ok bool
	// Scale normalization
	switch {
	case scale < minScale:
		coef, ok = coef.lsh(minScale - scale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		scale = minScale
	case scale > MaxScale:
		coef = coef.rshHalfEven(scale - MaxScale)
		scale = MaxScale
	}
	return newSafe(neg, coef, scale)
}

// newFromBint creates a new decimal from a *big.Int coefficient.
// This method uses overflowError to return descriptive errors.
func newFromBint(neg bool, coef *bint, scale, minScale int) (Decimal, error) {
	// Overflow validation
	prec := coef.prec()
	if prec-scale > MaxPrec-minScale {
		return Decimal{}, overflowError(prec, scale, minScale)
	}
	// Scale normalization
	switch {
	case scale < minScale:
		coef.lsh(coef, minScale-scale)
		scale = minScale
	case scale >= prec && scale > MaxScale: // no integer part
		coef.rshHalfEven(coef, scale-MaxScale)
		scale = MaxScale
	case prec > scale && prec > MaxPrec: // there is an integer part
		coef.rshHalfEven(coef, prec-MaxPrec)
		scale = MaxPrec - prec + scale
	}
	// Handling the rare case when rshHalfEven rounded
	// a 19-digit coefficient to a 20-digit coefficient.
	if coef.hasPrec(MaxPrec + 1) {
		return newFromBint(neg, coef, scale, minScale)
	}
	return newSafe(neg, coef.fint(), scale)
}

func overflowError(gotPrec, gotScale, wantScale int) error {
	maxDigits := MaxPrec - wantScale
	gotDigits := gotPrec - gotScale
	switch wantScale {
	case 0:
		return fmt.Errorf("%w: the integer part of a %T can have at most %v digits, but it has %v digits", errDecimalOverflow, Decimal{}, maxDigits, gotDigits)
	default:
		return fmt.Errorf("%w: with %v significant digits after the decimal point, the integer part of a %T can have at most %v digits, but it has %v digits", errDecimalOverflow, wantScale, Decimal{}, maxDigits, gotDigits)
	}
}

func unknownOverflowError() error {
	return fmt.Errorf("%w: the integer part of a %T can have at most %v digits, but it has significantly more digits", errDecimalOverflow, Decimal{}, MaxPrec)
}

// MustNew is like [New] but panics if the decimal cannot be constructed.
// It simplifies safe initialization of global variables holding decimals.
func MustNew(value int64, scale int) Decimal {
	d, err := New(value, scale)
	if err != nil {
		panic(fmt.Sprintf("New(%v, %v) failed: %v", value, scale, err))
	}
	return d
}

// New returns a decimal equal to value / 10^scale.
// New keeps trailing zeros in the fractional part to preserve scale.
//
// New returns an error if the scale is negative or greater than [MaxScale].
func New(value int64, scale int) (Decimal, error) {
	var coef fint
	var neg bool
	if value >= 0 {
		neg = false
		coef = fint(value)
	} else {
		neg = true
		if value == math.MinInt64 {
			coef = fint(math.MaxInt64) + 1
		} else {
			coef = fint(-value)
		}
	}
	return newSafe(neg, coef, scale)
}

// NewFromInt64 converts a pair of integers, representing the whole and
// fractional parts, to a (possibly rounded) decimal equal to whole + frac / 10^scale.
// NewFromInt64 removes all trailing zeros from the fractional part.
// This method is useful for converting amounts from [protobuf] format.
// See also method [Decimal.Int64].
//
// NewFromInt64 returns an error if:
//   - the whole and fractional parts have different signs;
//   - the scale is negative or greater than [MaxScale];
//   - frac / 10^scale is not within the range (-1, 1).
//
// [protobuf]: https://github.com/googleapis/googleapis/blob/master/google/type/money.proto
func NewFromInt64(whole, frac int64, scale int) (Decimal, error) {
	// Whole
	d, err := New(whole, 0)
	if err != nil {
		return Decimal{}, fmt.Errorf("converting integers: %w", err) // should never happen
	}
	// Fraction
	f, err := New(frac, scale)
	if err != nil {
		return Decimal{}, fmt.Errorf("converting integers: %w", err)
	}
	if !f.IsZero() {
		if !d.IsZero() && d.Sign() != f.Sign() {
			return Decimal{}, fmt.Errorf("converting integers: inconsistent signs")
		}
		if !f.WithinOne() {
			return Decimal{}, fmt.Errorf("converting integers: inconsistent fraction")
		}
		f = f.Trim(0)
		d, err = d.Add(f)
		if err != nil {
			return Decimal{}, fmt.Errorf("converting integers: %w", err) // should never happen
		}
	}
	return d, nil
}

// Int64 returns a pair of integers representing the whole and
// (possibly rounded) fractional parts of the decimal.
// If given scale is greater than the scale of the decimal, then the fractional part
// is zero-padded to the right.
// If given scale is smaller than the scale of the decimal, then the fractional part
// is rounded using [rounding half to even] (banker's rounding).
// The relationship between the decimal and the returned values can be expressed
// as d = whole + frac / 10^scale.
// This method is useful for converting amounts to [protobuf] format.
// See also constructor [NewFromInt64].
//
// If the result cannot be represented as a pair of int64 values,
// then false is returned.
//
// [rounding half to even]: https://en.wikipedia.org/wiki/Rounding#Rounding_half_to_even
// [protobuf]: https://github.com/googleapis/googleapis/blob/master/google/type/money.proto
func (d Decimal) Int64(scale int) (whole, frac int64, ok bool) {
	if scale < MinScale || scale > MaxScale {
		return 0, 0, false
	}
	x := d.coef
	y := pow10[d.Scale()]
	if scale < d.Scale() {
		x = x.rshHalfEven(d.Scale() - scale)
		y = pow10[scale]
	}
	q, r, ok := x.quoRem(y)
	if !ok {
		return 0, 0, false // Should never happen
	}
	if scale > d.Scale() {
		r, ok = r.lsh(scale - d.Scale())
		if !ok {
			return 0, 0, false // Should never happen
		}
	}
	if d.IsNeg() {
		if q > -math.MinInt64 || r > -math.MinInt64 {
			return 0, 0, false
		}
		//nolint:gosec
		return -int64(q), -int64(r), true
	}
	if q > math.MaxInt64 || r > math.MaxInt64 {
		return 0, 0, false
	}
	//nolint:gosec
	return int64(q), int64(r), true
}

// NewFromFloat64 converts a float to a (possibly rounded) decimal.
// See also method [Decimal.Float64].
//
// NewFromFloat64 returns an error if:
//   - the float is a special value (NaN or Inf);
//   - the integer part of the result has more than [MaxPrec] digits.
func NewFromFloat64(f float64) (Decimal, error) {
	// Float
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return Decimal{}, fmt.Errorf("converting float: special value %v", f)
	}
	text := make([]byte, 0, 32)
	text = strconv.AppendFloat(text, f, 'f', -1, 64)

	// Decimal
	d, err := parse(text)
	if err != nil {
		return Decimal{}, fmt.Errorf("converting float: %w", err)
	}
	return d, nil
}

// Float64 returns the nearest binary floating-point number rounded
// using [rounding half to even] (banker's rounding).
// See also constructor [NewFromFloat64].
//
// This conversion may lose data, as float64 has a smaller precision
// than the decimal type.
//
// [rounding half to even]: https://en.wikipedia.org/wiki/Rounding#Rounding_half_to_even
func (d Decimal) Float64() (f float64, ok bool) {
	s := d.String()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// MustParse is like [Parse] but panics if the string cannot be parsed.
// It simplifies safe initialization of global variables holding decimals.
func MustParse(s string) Decimal {
	d, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("Parse(%q) failed: %v", s, err))
	}
	return d
}

// Parse converts a string to a (possibly rounded) decimal.
// The input string must be in one of the following formats:
//
//	1.234
//	-1234
//	+0.000001234
//	1.83e5
//	0.22e-9
//
// The formal EBNF grammar for the supported format is as follows:
//
//	sign           ::= '+' | '-'
//	digits         ::= { '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' }
//	significand    ::= digits '.' digits | '.' digits | digits '.' | digits
//	exponent       ::= ('e' | 'E') [sign] digits
//	numeric-string ::= [sign] significand [exponent]
//
// Parse removes leading zeros from the integer part of the input string,
// but tries to maintain trailing zeros in the fractional part to preserve scale.
//
// Parse returns an error if:
//   - the string contains any whitespaces;
//   - the string is longer than 330 bytes;
//   - the exponent is less than -330 or greater than 330;
//   - the string does not represent a valid decimal number;
//   - the integer part of the result has more than [MaxPrec] digits.
func Parse(s string) (Decimal, error) {
	text := unsafe.Slice(unsafe.StringData(s), len(s))
	return parseExact(text, 0)
}

func parse(text []byte) (Decimal, error) {
	return parseExact(text, 0)
}

// ParseExact is similar to [Parse], but it allows you to specify how many digits
// after the decimal point should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for parsing monetary amounts, where the scale should be
// equal to or greater than the currency's scale.
func ParseExact(s string, scale int) (Decimal, error) {
	text := unsafe.Slice(unsafe.StringData(s), len(s))
	return parseExact(text, scale)
}

func parseExact(text []byte, scale int) (Decimal, error) {
	if len(text) > 330 {
		return Decimal{}, fmt.Errorf("parsing decimal: %w", errInvalidDecimal)
	}
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("parsing decimal: %w", errScaleRange)
	}
	d, err := parseFint(text, scale)
	if err != nil {
		d, err = parseBint(text, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("parsing decimal: %w", err)
		}
	}
	return d, nil
}

// parseFint parses a decimal string using uint64 arithmetic.
// parseFint does not support exponential notation to make it as fast as possible.
//
//nolint:gocyclo
func parseFint(text []byte, minScale int) (Decimal, error) {
	var pos int
	width := len(text)

	// Sign
	var neg bool
	switch {
	case pos == width:
		// skip
	case text[pos] == '-':
		neg = true
		pos++
	case text[pos] == '+':
		pos++
	}

	// Coefficient
	var coef fint
	var scale int
	var hasCoef, ok bool

	// Integer
	for pos < width && text[pos] >= '0' && text[pos] <= '9' {
		coef, ok = coef.fsa(1, text[pos]-'0')
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		pos++
		hasCoef = true
	}

	// Fraction
	if pos < width && text[pos] == '.' {
		pos++
		for pos < width && text[pos] >= '0' && text[pos] <= '9' {
			coef, ok = coef.fsa(1, text[pos]-'0')
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
			pos++
			scale++
			hasCoef = true
		}
	}

	if pos != width {
		return Decimal{}, fmt.Errorf("%w: unexpected character %q", errInvalidDecimal, text[pos])
	}
	if !hasCoef {
		return Decimal{}, fmt.Errorf("%w: no coefficient", errInvalidDecimal)
	}
	return newFromFint(neg, coef, scale, minScale)
}

// parseBint parses a decimal string using *big.Int arithmetic.
// parseBint supports exponential notation.
//
//nolint:gocyclo
func parseBint(text []byte, minScale int) (Decimal, error) {
	var pos int
	width := len(text)

	// Sign
	var neg bool
	switch {
	case pos == width:
		// skip
	case text[pos] == '-':
		neg = true
		pos++
	case text[pos] == '+':
		pos++
	}

	// Coefficient
	bcoef := getBint()
	defer putBint(bcoef)
	var fcoef fint
	var shift, scale int
	var hasCoef, ok bool

	bcoef.setFint(0)

	// Algorithm:
	// 	1. Add as many digits as possible to the uint64 coefficient (fast).
	// 	2. Once the uint64 coefficient has reached its maximum value,
	//     add it to the *big.Int coefficient (slow).
	// 	3. Repeat until all digits are processed.

	// Integer
	for pos < width && text[pos] >= '0' && text[pos] <= '9' {
		fcoef, ok = fcoef.fsa(1, text[pos]-'0')
		if !ok {
			return Decimal{}, errDecimalOverflow // Should never happen
		}
		pos++
		shift++
		hasCoef = true
		if fcoef.hasPrec(MaxPrec) {
			bcoef.fsa(bcoef, shift, fcoef)
			fcoef, shift = 0, 0
		}
	}

	// Fraction
	if pos < width && text[pos] == '.' {
		pos++
		for pos < width && text[pos] >= '0' && text[pos] <= '9' {
			fcoef, ok = fcoef.fsa(1, text[pos]-'0')
			if !ok {
				return Decimal{}, errDecimalOverflow // Should never happen
			}
			pos++
			scale++
			shift++
			hasCoef = true
			if fcoef.hasPrec(MaxPrec) {
				bcoef.fsa(bcoef, shift, fcoef)
				fcoef, shift = 0, 0
			}
		}
	}
	if shift > 0 {
		bcoef.fsa(bcoef, shift, fcoef)
	}

	// Exponent
	var exp int
	var eneg, hasExp, hasE bool
	if pos < width && (text[pos] == 'e' || text[pos] == 'E') {
		pos++
		hasE = true
		// Sign
		switch {
		case pos == width:
			// skip
		case text[pos] == '-':
			eneg = true
			pos++
		case text[pos] == '+':
			pos++
		}
		// Integer
		for pos < width && text[pos] >= '0' && text[pos] <= '9' {
			exp = exp*10 + int(text[pos]-'0')
			if exp > 330 {
				return Decimal{}, errInvalidDecimal
			}
			pos++
			hasExp = true
		}
	}

	if pos != width {
		return Decimal{}, fmt.Errorf("%w: unexpected character %q", errInvalidDecimal, text[pos])
	}
	if !hasCoef {
		return Decimal{}, fmt.Errorf("%w: no coefficient", errInvalidDecimal)
	}
	if hasE && !hasExp {
		return Decimal{}, fmt.Errorf("%w: no exponent", errInvalidDecimal)
	}

	if eneg {
		scale = scale + exp
	} else {
		scale = scale - exp
	}

	return newFromBint(neg, bcoef, scale, minScale)
}

// String implements the [fmt.Stringer] interface and returns
// a string representation of the decimal.
// The returned string does not use scientific or engineering notation and is
// formatted according to the following formal EBNF grammar:
//
//	sign           ::= '-'
//	digits         ::= { '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' }
//	significand    ::= digits '.' digits | digits
//	numeric-string ::= [sign] significand
//
// See also method [Decimal.Format].
//
// [fmt.Stringer]: https://pkg.go.dev/fmt#Stringer
func (d Decimal) String() string {
	return string(d.bytes())
}

// bytes returns a string representation of the decimal as a byte slice.
func (d Decimal) bytes() []byte {
	text := make([]byte, 0, 24)
	return d.append(text)
}

// append appends a string representation of the decimal to the byte slice.
func (d Decimal) append(text []byte) []byte {
	var buf [24]byte
	pos := len(buf) - 1
	coef := d.Coef()
	scale := d.Scale()

	// Coefficient
	for {
		buf[pos] = byte(coef%10) + '0'
		pos--
		coef /= 10
		if scale > 0 {
			scale--
			// Decimal point
			if scale == 0 {
				buf[pos] = '.'
				pos--
				// Leading 0
				if coef == 0 {
					buf[pos] = '0'
					pos--
				}
			}
		}
		if coef == 0 && scale == 0 {
			break
		}
	}

	// Sign
	if d.IsNeg() {
		buf[pos] = '-'
		pos--
	}

	return append(text, buf[pos+1:]...)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
// UnmarshalJSON supports the following types: [number] and [numeric string].
// See also constructor [Parse].
//
// [number]: https://datatracker.ietf.org/doc/html/rfc8259#section-6
// [numeric string]: https://datatracker.ietf.org/doc/html/rfc8259#section-7
// [json.Unmarshaler]: https://pkg.go.dev/encoding/json#Unmarshaler
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	var err error
	*d, err = parse(data)
	if err != nil {
		return fmt.Errorf("unmarshaling %T: %w", Decimal{}, err)
	}
	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
// MarshalJSON always returns a [numeric string].
// See also method [Decimal.String].
//
// [numeric string]: https://datatracker.ietf.org/doc/html/rfc8259#section-7
// [json.Marshaler]: https://pkg.go.dev/encoding/json#Marshaler
func (d Decimal) MarshalJSON() ([]byte, error) {
	text := make([]byte, 0, 26)
	text = append(text, '"')
	text = d.append(text)
	text = append(text, '"')
	return text, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// UnmarshalText supports only numeric strings.
// See also constructor [Parse].
//
// [encoding.TextUnmarshaler]: https://pkg.go.dev/encoding#TextUnmarshaler
func (d *Decimal) UnmarshalText(text []byte) error {
	var err error
	*d, err = parse(text)
	if err != nil {
		return fmt.Errorf("unmarshaling %T: %w", Decimal{}, err)
	}
	return nil
}

// AppendText implements the [encoding.TextAppender] interface.
// AppendText always appends a numeric string.
// See also method [Decimal.String].
//
// [encoding.TextAppender]: https://pkg.go.dev/encoding#TextAppender
func (d Decimal) AppendText(text []byte) ([]byte, error) {
	return d.append(text), nil
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// MarshalText always returns a numeric string.
// See also method [Decimal.String].
//
// [encoding.TextMarshaler]: https://pkg.go.dev/encoding#TextMarshaler
func (d Decimal) MarshalText() ([]byte, error) {
	return d.bytes(), nil
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
// UnmarshalBinary supports only numeric strings.
// See also constructor [Parse].
//
// [encoding.BinaryUnmarshaler]: https://pkg.go.dev/encoding#BinaryUnmarshaler
func (d *Decimal) UnmarshalBinary(data []byte) error {
	var err error
	*d, err = parse(data)
	if err != nil {
		return fmt.Errorf("unmarshaling %T: %w", Decimal{}, err)
	}
	return nil
}

// AppendBinary implements the [encoding.BinaryAppender] interface.
// AppendBinary always appends a numeric string.
// See also method [Decimal.String].
//
// [encoding.BinaryAppender]: https://pkg.go.dev/encoding#BinaryAppender
func (d Decimal) AppendBinary(data []byte) ([]byte, error) {
	return d.append(data), nil
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
// MarshalBinary always returns a numeric string.
// See also method [Decimal.String].
//
// [encoding.BinaryMarshaler]: https://pkg.go.dev/encoding#BinaryMarshaler
func (d Decimal) MarshalBinary() ([]byte, error) {
	return d.bytes(), nil
}

// UnmarshalBSONValue implements the [v2/bson.ValueUnmarshaler] interface.
// UnmarshalBSONValue supports the following [types]: Double, String, 32-bit Integer, 64-bit Integer, and [Decimal128].
//
// [v2/bson.ValueUnmarshaler]: https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/bson#ValueUnmarshaler
// [types]: https://bsonspec.org/spec.html
// [Decimal128]: https://github.com/mongodb/specifications/blob/master/source/bson-decimal128/decimal128.md
func (d *Decimal) UnmarshalBSONValue(typ byte, data []byte) error {
	// constants are from https://bsonspec.org/spec.html
	var err error
	switch typ {
	case 1:
		*d, err = parseBSONFloat64(data)
	case 2:
		*d, err = parseBSONString(data)
	case 10:
		// null, do nothing
	case 16:
		*d, err = parseBSONInt32(data)
	case 18:
		*d, err = parseBSONInt64(data)
	case 19:
		*d, err = parseIEEEDecimal128(data)
	default:
		err = fmt.Errorf("BSON type %d is not supported", typ)
	}
	if err != nil {
		err = fmt.Errorf("converting from BSON type %d to %T: %w", typ, Decimal{}, err)
	}
	return err
}

// MarshalBSONValue implements the [v2/bson.ValueMarshaler] interface.
// MarshalBSONValue always returns [Decimal128].
//
// [v2/bson.ValueMarshaler]: https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/bson#ValueMarshaler
// [Decimal128]: https://github.com/mongodb/specifications/blob/master/source/bson-decimal128/decimal128.md
func (d Decimal) MarshalBSONValue() (typ byte, data []byte, err error) {
	return 19, d.ieeeDecimal128(), nil
}

// parseBSONInt32 parses a BSON int32 to a decimal.
// The byte order of the input data must be little-endian.
func parseBSONInt32(data []byte) (Decimal, error) {
	if len(data) != 4 {
		return Decimal{}, fmt.Errorf("%w: invalid data length %v", errInvalidDecimal, len(data))
	}
	u := uint32(data[0])
	u |= uint32(data[1]) << 8
	u |= uint32(data[2]) << 16
	u |= uint32(data[3]) << 24
	i := int64(int32(u)) //nolint:gosec
	return New(i, 0)
}

// parseBSONInt64 parses a BSON int64 to a decimal.
// The byte order of the input data must be little-endian.
func parseBSONInt64(data []byte) (Decimal, error) {
	if len(data) != 8 {
		return Decimal{}, fmt.Errorf("%w: invalid data length %v", errInvalidDecimal, len(data))
	}
	u := uint64(data[0])
	u |= uint64(data[1]) << 8
	u |= uint64(data[2]) << 16
	u |= uint64(data[3]) << 24
	u |= uint64(data[4]) << 32
	u |= uint64(data[5]) << 40
	u |= uint64(data[6]) << 48
	u |= uint64(data[7]) << 56
	i := int64(u) //nolint:gosec
	return New(i, 0)
}

// parseBSONFloat64 parses a BSON float64 to a (possibly rounded) decimal.
// The byte order of the input data must be little-endian.
func parseBSONFloat64(data []byte) (Decimal, error) {
	if len(data) != 8 {
		return Decimal{}, fmt.Errorf("%w: invalid data length %v", errInvalidDecimal, len(data))
	}
	u := uint64(data[0])
	u |= uint64(data[1]) << 8
	u |= uint64(data[2]) << 16
	u |= uint64(data[3]) << 24
	u |= uint64(data[4]) << 32
	u |= uint64(data[5]) << 40
	u |= uint64(data[6]) << 48
	u |= uint64(data[7]) << 56
	f := math.Float64frombits(u)
	return NewFromFloat64(f)
}

// parseBSONString parses a BSON string to a (possibly rounded) decimal.
// The byte order of the input data must be little-endian.
func parseBSONString(data []byte) (Decimal, error) {
	if len(data) < 4 {
		return Decimal{}, fmt.Errorf("%w: invalid data length %v", errInvalidDecimal, len(data))
	}
	u := uint32(data[0])
	u |= uint32(data[1]) << 8
	u |= uint32(data[2]) << 16
	u |= uint32(data[3]) << 24
	l := int(int32(u)) //nolint:gosec
	if l < 1 || l > 330 || len(data) < l+4 {
		return Decimal{}, fmt.Errorf("%w: invalid string length %v", errInvalidDecimal, l)
	}
	if data[l+4-1] != 0 {
		return Decimal{}, fmt.Errorf("%w: invalid null terminator %v", errInvalidDecimal, data[l+4-1])
	}
	s := string(data[4 : l+4-1])
	return Parse(s)
}

// parseIEEEDecimal128 converts a 128-bit IEEE 754-2008 decimal
// floating point with binary integer decimal encoding to
// a (possibly rounded) decimal.
// The byte order of the input data must be little-endian.
//
// parseIEEEDecimal128 returns an error if:
//   - the data length is not equal to 16 bytes;
//   - the decimal a special value (NaN or Inf);
//   - the integer part of the result has more than [MaxPrec] digits.
func parseIEEEDecimal128(data []byte) (Decimal, error) {
	if len(data) != 16 {
		return Decimal{}, fmt.Errorf("%w: invalid data length %v", errInvalidDecimal, len(data))
	}
	if data[15]&0b0111_1100 == 0b0111_1100 {
		return Decimal{}, fmt.Errorf("%w: special value NaN", errInvalidDecimal)
	}
	if data[15]&0b0111_1100 == 0b0111_1000 {
		return Decimal{}, fmt.Errorf("%w: special value Inf", errInvalidDecimal)
	}
	if data[15]&0b0110_0000 == 0b0110_0000 {
		return Decimal{}, fmt.Errorf("%w: unsupported encoding", errInvalidDecimal)
	}

	// Sign
	neg := data[15]&0b1000_0000 == 0b1000_0000

	// Scale
	var scale int
	scale |= int(data[14]) >> 1
	scale |= int(data[15]&0b0111_1111) << 7
	scale = 6176 - scale

	// TODO fint optimization

	// Coefficient
	coef := getBint()
	defer putBint(coef)

	buf := make([]byte, 15)
	for i := range 15 {
		buf[i] = data[14-i]
	}
	buf[0] &= 0b0000_0001
	coef.setBytes(buf)

	// Scale normalization
	if coef.sign() == 0 {
		scale = max(scale, MinScale)
	}

	return newFromBint(neg, coef, scale, 0)
}

// ieeeDecimal128 returns a 128-bit IEEE 754-2008 decimal
// floating point with binary integer decimal encoding.
// The byte order of the result is little-endian.
func (d Decimal) ieeeDecimal128() []byte {
	var buf [16]byte
	scale := d.Scale()
	coef := d.Coef()

	// Sign
	if d.IsNeg() {
		buf[15] = 0b1000_0000
	}

	// Scale
	scale = 6176 - scale
	buf[15] |= byte((scale >> 7) & 0b0111_1111)
	buf[14] |= byte((scale << 1) & 0b1111_1110)

	// Coefficient
	for i := range 8 {
		buf[i] = byte(coef & 0b1111_1111)
		coef >>= 8
	}

	return buf[:]
}

// Scan implements the [sql.Scanner] interface.
//
// [sql.Scanner]: https://pkg.go.dev/database/sql#Scanner
func (d *Decimal) Scan(value any) error {
	var err error
	switch value := value.(type) {
	case string:
		*d, err = Parse(value)
	case int64:
		*d, err = New(value, 0)
	case float64:
		*d, err = NewFromFloat64(value)
	case []byte:
		// Special case: MySQL driver sends DECIMAL as []byte
		*d, err = parse(value)
	case float32:
		// Special case: MySQL driver sends FLOAT as float32
		*d, err = NewFromFloat64(float64(value))
	case uint64:
		// Special case: ClickHouse driver sends 0 as uint64
		*d, err = newSafe(false, fint(value), 0)
	case nil:
		err = fmt.Errorf("%T does not support null values, use %T or *%T", Decimal{}, NullDecimal{}, Decimal{})
	default:
		err = fmt.Errorf("type %T is not supported", value)
	}
	if err != nil {
		err = fmt.Errorf("converting from %T to %T: %w", value, Decimal{}, err)
	}
	return err
}

// Value implements the [driver.Valuer] interface.
//
// [driver.Valuer]: https://pkg.go.dev/database/sql/driver#Valuer
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Format implements the [fmt.Formatter] interface.
// The following [format verbs] are available:
//
//	| Verb       | Example | Description    |
//	| ---------- | ------- | -------------- |
//	| %f, %s, %v | 5.67    | Decimal        |
//	| %q         | "5.67"  | Quoted decimal |
//	| %k         | 567%    | Percentage     |
//
// The following format flags can be used with all verbs: '+', ' ', '0', '-'.
//
// Precision is only supported for %f and %k verbs.
// For %f verb, the default precision is equal to the actual scale of the decimal,
// whereas, for verb %k the default precision is the actual scale of the decimal minus 2.
//
// [format verbs]: https://pkg.go.dev/fmt#hdr-Printing
// [fmt.Formatter]: https://pkg.go.dev/fmt#Formatter
//
//nolint:gocyclo
func (d Decimal) Format(state fmt.State, verb rune) {
	var err error

	// Percentage multiplier
	if verb == 'k' || verb == 'K' {
		d, err = d.Mul(Hundred)
		if err != nil {
			// This panic is handled inside the fmt package.
			panic(fmt.Errorf("formatting percent: %w", err))
		}
	}

	// Rescaling
	var tzeros int
	if verb == 'f' || verb == 'F' || verb == 'k' || verb == 'K' {
		var scale int
		switch p, ok := state.Precision(); {
		case ok:
			scale = p
		case verb == 'k' || verb == 'K':
			scale = d.Scale() - 2
		case verb == 'f' || verb == 'F':
			scale = d.Scale()
		}
		scale = max(scale, MinScale)
		switch {
		case scale < d.Scale():
			d = d.Round(scale)
		case scale > d.Scale():
			tzeros = scale - d.Scale()
		}
	}

	// Integer and fractional digits
	var intdigs int
	fracdigs := d.Scale()
	if dprec := d.Prec(); dprec > fracdigs {
		intdigs = dprec - fracdigs
	}
	if d.WithinOne() {
		intdigs++ // leading 0
	}

	// Decimal point
	var dpoint int
	if fracdigs > 0 || tzeros > 0 {
		dpoint = 1
	}

	// Arithmetic sign
	var rsign int
	if d.IsNeg() || state.Flag('+') || state.Flag(' ') {
		rsign = 1
	}

	// Percentage sign
	var psign int
	if verb == 'k' || verb == 'K' {
		psign = 1
	}

	// Openning and closing quotes
	var lquote, tquote int
	if verb == 'q' || verb == 'Q' {
		lquote, tquote = 1, 1
	}

	// Calculating padding
	width := lquote + rsign + intdigs + dpoint + fracdigs + tzeros + psign + tquote
	var lspaces, tspaces, lzeros int
	if w, ok := state.Width(); ok && w > width {
		switch {
		case state.Flag('-'):
			tspaces = w - width
		case state.Flag('0'):
			lzeros = w - width
		default:
			lspaces = w - width
		}
		width = w
	}

	buf := make([]byte, width)
	pos := width - 1

	// Trailing spaces
	for range tspaces {
		buf[pos] = ' '
		pos--
	}

	// Closing quote
	for range tquote {
		buf[pos] = '"'
		pos--
	}

	// Percentage sign
	for range psign {
		buf[pos] = '%'
		pos--
	}

	// Trailing zeros
	for range tzeros {
		buf[pos] = '0'
		pos--
	}

	// Fractional digits
	dcoef := d.Coef()
	for range fracdigs {
		buf[pos] = byte(dcoef%10) + '0'
		pos--
		dcoef /= 10
	}

	// Decimal point
	for range dpoint {
		buf[pos] = '.'
		pos--
	}

	// Integer digits
	for range intdigs {
		buf[pos] = byte(dcoef%10) + '0'
		pos--
		dcoef /= 10
	}

	// Leading zeros
	for range lzeros {
		buf[pos] = '0'
		pos--
	}

	// Arithmetic sign
	for range rsign {
		if d.IsNeg() {
			buf[pos] = '-'
		} else if state.Flag(' ') {
			buf[pos] = ' '
		} else {
			buf[pos] = '+'
		}
		pos--
	}

	// Opening quote
	for range lquote {
		buf[pos] = '"'
		pos--
	}

	// Leading spaces
	for range lspaces {
		buf[pos] = ' '
		pos--
	}

	// Writing result
	//nolint:errcheck
	switch verb {
	case 'q', 'Q', 's', 'S', 'v', 'V', 'f', 'F', 'k', 'K':
		state.Write(buf)
	default:
		state.Write([]byte("%!"))
		state.Write([]byte{byte(verb)})
		state.Write([]byte("(decimal.Decimal="))
		state.Write(buf)
		state.Write([]byte(")"))
	}
}

// Zero returns a decimal with a value of 0, having the same scale as decimal d.
// See also methods [Decimal.One], [Decimal.ULP].
func (d Decimal) Zero() Decimal {
	return newUnsafe(false, 0, d.Scale())
}

// One returns a decimal with a value of 1, having the same scale as decimal d.
// See also methods [Decimal.Zero], [Decimal.ULP].
func (d Decimal) One() Decimal {
	return newUnsafe(false, pow10[d.Scale()], d.Scale())
}

// ULP (Unit in the Last Place) returns the smallest representable positive
// difference between two decimals with the same scale as decimal d.
// It can be useful for implementing rounding and comparison algorithms.
// See also methods [Decimal.Zero], [Decimal.One].
func (d Decimal) ULP() Decimal {
	return newUnsafe(false, 1, d.Scale())
}

// Prec returns the number of digits in the coefficient.
// See also method [Decimal.Coef].
func (d Decimal) Prec() int {
	return d.coef.prec()
}

// Coef returns the coefficient of the decimal.
// See also method [Decimal.Prec].
func (d Decimal) Coef() uint64 {
	return uint64(d.coef)
}

// Scale returns the number of digits after the decimal point.
// See also methods [Decimal.Prec], [Decimal.MinScale].
func (d Decimal) Scale() int {
	return int(d.scale)
}

// MinScale returns the smallest scale that the decimal can be rescaled to
// without rounding.
// See also method [Decimal.Trim].
func (d Decimal) MinScale() int {
	// Special case: zero
	if d.IsZero() {
		return MinScale
	}
	// General case
	dcoef := d.coef
	return max(MinScale, d.Scale()-dcoef.ntz())
}

// IsInt returns true if there are no significant digits after the decimal point.
func (d Decimal) IsInt() bool {
	return d.Scale() == 0 || d.coef%pow10[d.Scale()] == 0
}

// IsOne returns:
//
//	true  if d = -1 or d = 1
//	false otherwise
func (d Decimal) IsOne() bool {
	return d.coef == pow10[d.Scale()]
}

// WithinOne returns:
//
//	true  if -1 < d < 1
//	false otherwise
func (d Decimal) WithinOne() bool {
	return d.coef < pow10[d.Scale()]
}

// Round returns a decimal rounded to the specified number of digits after
// the decimal point using [rounding half to even] (banker's rounding).
// If the given scale is negative, it is redefined to zero.
// For financial calculations, the scale should be equal to or greater than
// the scale of the currency.
// See also method [Decimal.Rescale].
//
// [rounding half to even]: https://en.wikipedia.org/wiki/Rounding#Rounding_half_to_even
func (d Decimal) Round(scale int) Decimal {
	scale = max(scale, MinScale)
	if scale >= d.Scale() {
		return d
	}
	coef := d.coef
	coef = coef.rshHalfEven(d.Scale() - scale)
	return newUnsafe(d.IsNeg(), coef, scale)
}

// Pad returns a decimal zero-padded to the specified number of digits after
// the decimal point.
// The total number of digits in the result is limited by [MaxPrec].
// See also method [Decimal.Trim].
func (d Decimal) Pad(scale int) Decimal {
	scale = min(scale, MaxScale, MaxPrec-d.Prec()+d.Scale())
	if scale <= d.Scale() {
		return d
	}
	coef := d.coef
	coef, ok := coef.lsh(scale - d.Scale())
	if !ok {
		return d // Should never happen
	}
	return newUnsafe(d.IsNeg(), coef, scale)
}

// Rescale returns a decimal rounded or zero-padded to the given number of digits
// after the decimal point.
// If the given scale is negative, it is redefined to zero.
// For financial calculations, the scale should be equal to or greater than
// the scale of the currency.
// See also methods [Decimal.Round], [Decimal.Pad].
func (d Decimal) Rescale(scale int) Decimal {
	if scale > d.Scale() {
		return d.Pad(scale)
	}
	return d.Round(scale)
}

// Quantize returns a decimal rescaled to the same scale as decimal e.
// The sign and the coefficient of decimal e are ignored.
// See also methods [Decimal.SameScale] and [Decimal.Rescale].
func (d Decimal) Quantize(e Decimal) Decimal {
	return d.Rescale(e.Scale())
}

// SameScale returns true if decimals have the same scale.
// See also methods [Decimal.Scale], [Decimal.Quantize].
func (d Decimal) SameScale(e Decimal) bool {
	return d.Scale() == e.Scale()
}

// Trunc returns a decimal truncated to the specified number of digits
// after the decimal point using [rounding toward zero].
// If the given scale is negative, it is redefined to zero.
// For financial calculations, the scale should be equal to or greater than
// the scale of the currency.
//
// [rounding toward zero]: https://en.wikipedia.org/wiki/Rounding#Rounding_toward_zero
func (d Decimal) Trunc(scale int) Decimal {
	scale = max(scale, MinScale)
	if scale >= d.Scale() {
		return d
	}
	coef := d.coef
	coef = coef.rshDown(d.Scale() - scale)
	return newUnsafe(d.IsNeg(), coef, scale)
}

// Trim returns a decimal with trailing zeros removed up to the given number of
// digits after the decimal point.
// If the given scale is negative, it is redefined to zero.
// See also method [Decimal.Pad].
func (d Decimal) Trim(scale int) Decimal {
	if d.Scale() <= scale {
		return d
	}
	scale = max(scale, d.MinScale())
	return d.Trunc(scale)
}

// Ceil returns a decimal rounded up to the given number of digits
// after the decimal point using [rounding toward positive infinity].
// If the given scale is negative, it is redefined to zero.
// For financial calculations, the scale should be equal to or greater than
// the scale of the currency.
// See also method [Decimal.Floor].
//
// [rounding toward positive infinity]: https://en.wikipedia.org/wiki/Rounding#Rounding_up
func (d Decimal) Ceil(scale int) Decimal {
	scale = max(scale, MinScale)
	if scale >= d.Scale() {
		return d
	}
	coef := d.coef
	if d.IsNeg() {
		coef = coef.rshDown(d.Scale() - scale)
	} else {
		coef = coef.rshUp(d.Scale() - scale)
	}
	return newUnsafe(d.IsNeg(), coef, scale)
}

// Floor returns a decimal rounded down to the specified number of digits
// after the decimal point using [rounding toward negative infinity].
// If the given scale is negative, it is redefined to zero.
// For financial calculations, the scale should be equal to or greater than
// the scale of the currency.
// See also method [Decimal.Ceil].
//
// [rounding toward negative infinity]: https://en.wikipedia.org/wiki/Rounding#Rounding_down
func (d Decimal) Floor(scale int) Decimal {
	scale = max(scale, MinScale)
	if scale >= d.Scale() {
		return d
	}
	coef := d.coef
	if d.IsNeg() {
		coef = coef.rshUp(d.Scale() - scale)
	} else {
		coef = coef.rshDown(d.Scale() - scale)
	}
	return newUnsafe(d.IsNeg(), coef, scale)
}

// Neg returns a decimal with the opposite sign.
func (d Decimal) Neg() Decimal {
	return newUnsafe(!d.IsNeg(), d.coef, d.Scale())
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	return newUnsafe(false, d.coef, d.Scale())
}

// CopySign returns a decimal with the same sign as decimal e.
// CopySign treates 0 as positive.
// See also method [Decimal.Sign].
func (d Decimal) CopySign(e Decimal) Decimal {
	if d.IsNeg() == e.IsNeg() {
		return d
	}
	return d.Neg()
}

// Sign returns:
//
//	-1 if d < 0
//	 0 if d = 0
//	+1 if d > 0
//
// See also methods [Decimal.IsPos], [Decimal.IsNeg], [Decimal.IsZero].
func (d Decimal) Sign() int {
	switch {
	case d.neg:
		return -1
	case d.coef == 0:
		return 0
	}
	return 1
}

// IsPos returns:
//
//	true  if d > 0
//	false otherwise
func (d Decimal) IsPos() bool {
	return d.coef != 0 && !d.neg
}

// IsNeg returns:
//
//	true  if d < 0
//	false otherwise
func (d Decimal) IsNeg() bool {
	return d.neg
}

// IsZero returns:
//
//	true  if d = 0
//	false otherwise
func (d Decimal) IsZero() bool {
	return d.coef == 0
}

// Prod returns the (possibly rounded) product of decimals.
// It computes d1 * d2 * ... * dn with at least double precision
// during the intermediate rounding.
//
// Prod returns an error if:
//   - no arguments are provided;
//   - the integer part of the result has more than [MaxPrec] digits.
func Prod(d ...Decimal) (Decimal, error) {
	// Special cases
	switch len(d) {
	case 0:
		return Decimal{}, fmt.Errorf("computing [prod([])]: %w", errInvalidOperation)
	case 1:
		return d[0], nil
	}

	// General case
	e, err := prodFint(d...)
	if err != nil {
		e, err = prodBint(d...)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [prod(%v)]: %w", d, err)
		}
	}

	return e, nil
}

// prodFint computes the product of decimals using uint64 arithmetic.
func prodFint(d ...Decimal) (Decimal, error) {
	ecoef := One.coef
	escale := One.Scale()
	eneg := One.IsNeg()

	for _, f := range d {
		fcoef := f.coef

		// Compute e = e * f
		var ok bool
		ecoef, ok = ecoef.mul(fcoef)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		eneg = eneg != f.IsNeg()
		escale = escale + f.Scale()
	}

	return newFromFint(eneg, ecoef, escale, 0)
}

// prodBint computes the product of decimals using *big.Int arithmetic.
func prodBint(d ...Decimal) (Decimal, error) {
	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	ecoef.setFint(One.coef)
	escale := One.Scale()
	eneg := One.IsNeg()

	for _, f := range d {
		fcoef.setFint(f.coef)

		// Compute e = e * f
		ecoef.mul(ecoef, fcoef)
		eneg = eneg != f.IsNeg()
		escale = escale + f.Scale()

		// Intermediate truncation
		if escale > bscale {
			ecoef.rshDown(ecoef, escale-bscale)
			escale = bscale
		}

		// Check if e >= 10^59
		if ecoef.hasPrec(len(bpow10)) {
			return Decimal{}, unknownOverflowError()
		}
	}

	return newFromBint(eneg, ecoef, escale, 0)
}

// Mean returns the (possibly rounded) mean of decimals.
// It computes (d1 + d2 + ... + dn) / n with at least double precision
// during the intermediate rounding.
//
// Mean returns an error if:
//   - no arguments are provided;
//   - the integer part of the result has more than [MaxPrec] digits.
func Mean(d ...Decimal) (Decimal, error) {
	// Special cases
	switch len(d) {
	case 0:
		return Decimal{}, fmt.Errorf("computing [mean([])]: %w", errInvalidOperation)
	case 1:
		return d[0], nil
	}

	// General case
	e, err := meanFint(d...)
	if err != nil {
		e, err = meanBint(d...)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [mean(%v)]: %w", d, err)
		}
	}

	// Preferred scale
	scale := 0
	for _, f := range d {
		scale = max(scale, f.Scale())
	}
	e = e.Trim(scale)

	return e, nil
}

// meanFint computes the mean of decimals using uint64 arithmetic.
func meanFint(d ...Decimal) (Decimal, error) {
	ecoef := Zero.coef
	escale := Zero.Scale()
	eneg := Zero.IsNeg()

	ncoef := fint(len(d))

	for _, f := range d {
		fcoef := f.coef

		// Alignment
		var ok bool
		switch {
		case escale > f.Scale():
			fcoef, ok = fcoef.lsh(escale - f.Scale())
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
		case escale < f.Scale():
			ecoef, ok = ecoef.lsh(f.Scale() - escale)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
			escale = f.Scale()
		}

		// Compute e = e + f
		if eneg == f.IsNeg() {
			ecoef, ok = ecoef.add(fcoef)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
		} else {
			if fcoef > ecoef {
				eneg = f.IsNeg()
			}
			ecoef = ecoef.subAbs(fcoef)
		}
	}

	// Alignment
	var ok bool
	if shift := MaxPrec - ecoef.prec(); shift > 0 {
		ecoef, ok = ecoef.lsh(shift)
		if !ok {
			return Decimal{}, errDecimalOverflow // Should never happen
		}
		escale = escale + shift
	}

	// Compute e = e / n
	ecoef, ok = ecoef.quo(ncoef)
	if !ok {
		return Decimal{}, errInexactDivision
	}

	return newFromFint(eneg, ecoef, escale, 0)
}

// meanBint computes the mean of decimals using *big.Int arithmetic.
func meanBint(d ...Decimal) (Decimal, error) {
	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	ncoef := getBint()
	defer putBint(ncoef)

	ecoef.setFint(Zero.coef)
	escale := Zero.Scale()
	eneg := Zero.IsNeg()
	ncoef.setInt64(int64(len(d)))

	for _, f := range d {
		fcoef.setFint(f.coef)

		// Alignment
		switch {
		case escale > f.Scale():
			fcoef.lsh(fcoef, escale-f.Scale())
		case escale < f.Scale():
			ecoef.lsh(ecoef, f.Scale()-escale)
			escale = f.Scale()
		}

		// Compute e = e + f
		if eneg == f.IsNeg() {
			ecoef.add(ecoef, fcoef)
		} else {
			if fcoef.cmp(ecoef) > 0 {
				eneg = f.IsNeg()
			}
			ecoef.subAbs(ecoef, fcoef)
		}
	}

	// Alignment
	ecoef.lsh(ecoef, bscale-escale)

	// Compute e = e / n
	ecoef.quo(ecoef, ncoef)

	return newFromBint(eneg, ecoef, bscale, 0)
}

// Mul returns the (possibly rounded) product of decimals d and e.
//
// Mul returns an overflow error if the integer part of the result has
// more than [MaxPrec] digits.
func (d Decimal) Mul(e Decimal) (Decimal, error) {
	return d.MulExact(e, 0)
}

// MulExact is similar to [Decimal.Mul], but it allows you to specify the number
// of digits after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will
// return an overflow error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) MulExact(e Decimal, scale int) (Decimal, error) {
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("computing [%v * %v]: %w", d, e, errScaleRange)
	}

	// General case
	f, err := d.mulFint(e, scale)
	if err != nil {
		f, err = d.mulBint(e, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v * %v]: %w", d, e, err)
		}
	}
	return f, nil
}

// mulFint computes the product of two decimals using uint64 arithmetic.
func (d Decimal) mulFint(e Decimal, minScale int) (Decimal, error) {
	dcoef := d.coef
	dscale := d.Scale()
	dneg := d.IsNeg()

	ecoef := e.coef

	// Compute d = d * e
	dcoef, ok := dcoef.mul(ecoef)
	if !ok {
		return Decimal{}, errDecimalOverflow
	}
	dscale = dscale + e.Scale()
	dneg = dneg != e.IsNeg()

	return newFromFint(dneg, dcoef, dscale, minScale)
}

// mulBint computes the product of two decimals using *big.Int arithmetic.
func (d Decimal) mulBint(e Decimal, minScale int) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	dscale := d.Scale()
	dneg := d.IsNeg()
	ecoef.setFint(e.coef)

	// Compute d = d * e
	dcoef.mul(dcoef, ecoef)
	dneg = dneg != e.IsNeg()
	dscale = dscale + e.Scale()

	return newFromBint(dneg, dcoef, dscale, minScale)
}

// Pow returns the (possibly rounded) decimal raised to the given decimal power.
// If zero is raised to zero power then the result is one.
//
// Pow returns an error if:
//   - the integer part of the result has more than [MaxPrec] digits;
//   - zero is raised to a negative power;
//   - negative is raised to a fractional power.
func (d Decimal) Pow(e Decimal) (Decimal, error) {
	// Special case: zero to a negative power
	if e.IsNeg() && d.IsZero() {
		return Decimal{}, fmt.Errorf("computing [%v^%v]: %w: zero to negative power", d, e, errInvalidOperation)
	}

	// Special case: integer power
	if e.IsInt() {
		power := e.Trunc(0).Coef()
		f, err := d.powIntFint(power, e.IsNeg())
		if err != nil {
			f, err = d.powIntBint(power, e.IsNeg())
			if err != nil {
				return Decimal{}, fmt.Errorf("computing [%v^%v]: %w", d, e, err)
			}
		}

		// Preferred scale
		if e.IsNeg() {
			f = f.Trim(0)
		}

		return f, nil
	}

	// Special case: zero to a fractional power
	if d.IsZero() {
		return newSafe(false, 0, 0)
	}

	// Special case: negative to a fractional power
	if d.IsNeg() {
		return Decimal{}, fmt.Errorf("computing [%v^%v]: %w: negative to fractional power", d, e, errInvalidOperation)
	}

	// General case
	f, err := d.powBint(e)
	if err != nil {
		return Decimal{}, fmt.Errorf("computing [%v^%v]: %w", d, e, err)
	}

	return f, nil
}

// powBint computes the power of a decimal using *big.Int arithmetic.
func (d Decimal) powBint(e Decimal) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	dcoef.setFint(d.coef)
	ecoef.setFint(e.coef)
	inv := false

	// Alignment
	if d.WithinOne() {
		// Compute d = ⌊1 / d⌋
		dcoef.quo(bpow10[bscale+d.Scale()], dcoef)
		inv = true
	} else {
		dcoef.lsh(dcoef, bscale-d.Scale())
	}

	// Compute f = log(d)
	fcoef.log(dcoef)

	// Compute f = ⌊f * e⌋
	fcoef.mul(fcoef, ecoef)
	fcoef.rshDown(fcoef, e.Scale())
	inv = inv != e.IsNeg()

	// Check if f <= -100 or f >= 100
	if fcoef.hasPrec(3 + bscale) {
		if !inv {
			return Decimal{}, unknownOverflowError()
		}
		return newSafe(false, 0, MaxScale)
	}

	// Compute f = exp(f)
	fcoef.exp(fcoef)

	if inv {
		// Compute f = ⌊1 / f⌋
		fcoef.quo(bpow10[2*bscale], fcoef)
	}

	return newFromBint(false, fcoef, bscale, 0)
}

// PowInt returns the (possibly rounded) decimal raised to the given integer power.
// If zero is raised to zero power then the result is one.
//
// PowInt returns an error if:
//   - the integer part of the result has more than [MaxPrec] digits;
//   - zero is raised to a negative power.
func (d Decimal) PowInt(power int) (Decimal, error) {
	var pow uint64
	var neg bool
	if power >= 0 {
		neg = false
		pow = uint64(power)
	} else {
		neg = true
		if power == math.MinInt {
			pow = uint64(math.MaxInt) + 1
		} else {
			pow = uint64(-power)
		}
	}

	// Special case: zero to a negative power
	if neg && d.IsZero() {
		return Decimal{}, fmt.Errorf("computing [%v^%v]: %w: zero to negative power", d, power, errInvalidOperation)
	}

	// General case
	e, err := d.powIntFint(pow, neg)
	if err != nil {
		e, err = d.powIntBint(pow, neg)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v^%v]: %w", d, power, err)
		}
	}

	// Preferred scale
	if neg {
		e = e.Trim(0)
	}

	return e, nil
}

// powIntFint computes the integer power of a decimal using uint64 arithmetic.
// powIntFint does not support negative powers.
func (d Decimal) powIntFint(pow uint64, inv bool) (Decimal, error) {
	if inv {
		return Decimal{}, errInvalidOperation
	}

	dcoef := d.coef
	dneg := d.IsNeg()
	dscale := d.Scale()

	ecoef := One.coef
	eneg := One.IsNeg()
	escale := One.Scale()

	// Exponentiation by squaring
	var ok bool
	for pow > 0 {
		if pow%2 == 1 {
			pow = pow - 1

			// Compute e = e * d
			ecoef, ok = ecoef.mul(dcoef)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
			eneg = eneg != dneg
			escale = escale + dscale
		}
		if pow > 0 {
			pow = pow / 2

			// Compute d = d * d
			dcoef, ok = dcoef.mul(dcoef)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
			dneg = false
			dscale = dscale * 2
		}
	}

	return newFromFint(eneg, ecoef, escale, 0)
}

// powIntBint computes the integer power of a decimal using *big.Int arithmetic.
// powIntBint supports negative powers.
func (d Decimal) powIntBint(pow uint64, inv bool) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	dneg := d.IsNeg()
	dscale := d.Scale()

	ecoef.setFint(One.coef)
	eneg := One.IsNeg()
	escale := One.Scale()

	// Exponentiation by squaring
	for pow > 0 {
		if pow%2 == 1 {
			pow = pow - 1

			// Compute e = e * d
			ecoef.mul(ecoef, dcoef)
			eneg = eneg != dneg
			escale = escale + dscale

			// Intermediate truncation
			if escale > bscale {
				ecoef.rshDown(ecoef, escale-bscale)
				escale = bscale
			}

			// Check if e <= -10^59 or e >= 10^59
			if ecoef.hasPrec(len(bpow10)) {
				if !inv {
					return Decimal{}, unknownOverflowError()
				}
				return newSafe(false, 0, MaxScale)
			}
		}
		if pow > 0 {
			pow = pow / 2

			// Compute d = d * d
			dcoef.mul(dcoef, dcoef)
			dneg = false
			dscale = dscale * 2

			// Intermediate truncation
			if dscale > bscale {
				dcoef.rshDown(dcoef, dscale-bscale)
				dscale = bscale
			}

			// Check if d <= -10^59 or d >= 10^59
			if dcoef.hasPrec(len(bpow10)) {
				if !inv {
					return Decimal{}, unknownOverflowError()
				}
				return newSafe(false, 0, MaxScale)
			}
		}
	}

	if inv {
		if ecoef.sign() == 0 {
			return Decimal{}, unknownOverflowError()
		}

		// Compute e = ⌊1 / e⌋
		ecoef.quo(bpow10[bscale+escale], ecoef)
		escale = bscale
	}

	return newFromBint(eneg, ecoef, escale, 0)
}

// Sqrt computes the (possibly rounded) square root of a decimal.
// d.Sqrt() is significantly faster than d.Pow(0.5).
//
// Sqrt returns an error if the decimal is negative.
func (d Decimal) Sqrt() (Decimal, error) {
	// Special case: negative
	if d.IsNeg() {
		return Decimal{}, fmt.Errorf("computing sqrt(%v): %w: square root of negative", d, errInvalidOperation)
	}

	// Special case: zero
	if d.IsZero() {
		return newSafe(false, 0, d.Scale()/2)
	}

	// General case
	e, err := d.sqrtBint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing sqrt(%v): %w", d, err)
	}

	// Preferred scale
	e = e.Trim(d.Scale() / 2)

	return e, nil
}

// sqrtBint computes the square root of a decimal using *big.Int arithmetic.
func (d Decimal) sqrtBint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	dcoef.setFint(d.coef)
	fcoef.setFint(0)

	// Alignment
	dcoef.lsh(dcoef, 2*bscale-d.Scale())

	// Initial guess is calculated as 10^(n/2),
	// where n is the position of the most significant digit.
	n := dcoef.prec() - 2*bscale
	ecoef.setBint(bpow10[n/2+bscale])

	// Newton's method
	for range 50 {
		if ecoef.cmp(fcoef) == 0 {
			break
		}
		fcoef.setBint(ecoef)
		ecoef.quo(dcoef, ecoef)
		ecoef.add(ecoef, fcoef)
		ecoef.hlf(ecoef)
	}

	return newFromBint(false, ecoef, bscale, 0)
}

// Log2 returns the (possibly rounded) binary logarithm of a decimal.
//
// Log2 returns an error if the decimal is zero or negative.
func (d Decimal) Log2() (Decimal, error) {
	// Special case: zero or negative
	if !d.IsPos() {
		return Decimal{}, fmt.Errorf("computing log2(%v): %w: logarithm of non-positive", d, errInvalidOperation)
	}

	// Special case: one
	if d.IsOne() {
		return newSafe(false, 0, 0)
	}

	// General case
	e, err := d.log2Bint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing log2(%v): %w", d, err)
	}

	// Preferred scale
	if e.IsInt() {
		// According to the GDA, only integer powers of 2 should be trimmed to zero scale.
		// However, such validation is slow, so we will trim all integers.
		e = e.Trunc(0)
	}

	return e, nil
}

// log2Bint computes the binary logarithm of a decimal using *big.Int arithmetic.
func (d Decimal) log2Bint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	eneg := false

	// Alignment
	if d.WithinOne() {
		// Compute d = ⌊1 / d⌋
		dcoef.quo(bpow10[bscale+d.Scale()], dcoef)
		eneg = true
	} else {
		dcoef.lsh(dcoef, bscale-d.Scale())
	}

	// Compute e = log(d)
	ecoef.log(dcoef)

	// Compute e = e / log(2)
	ecoef.lsh(ecoef, bscale)
	ecoef.quo(ecoef, blog[2])

	return newFromBint(eneg, ecoef, bscale, 0)
}

// Log10 returns the (possibly rounded) decimal logarithm of a decimal.
//
// Log10 returns an error if the decimal is zero or negative.
func (d Decimal) Log10() (Decimal, error) {
	// Special case: zero or negative
	if !d.IsPos() {
		return Decimal{}, fmt.Errorf("computing log10(%v): %w: logarithm of non-positive", d, errInvalidOperation)
	}

	// Special case: one
	if d.IsOne() {
		return newSafe(false, 0, 0)
	}

	// General case
	e, err := d.log10Bint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing log10(%v): %w", d, err)
	}

	// Preferred scale
	if e.IsInt() {
		// According to the GDA, only integer powers of 10 should be trimmed to zero scale.
		// However, such validation is slow, so we will trim all integers.
		e = e.Trunc(0)
	}

	return e, nil
}

// log10Bint computes the decimal logarithm of a decimal using *big.Int arithmetic.
func (d Decimal) log10Bint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	eneg := false

	// Alignment
	if d.WithinOne() {
		// Compute d = ⌊1 / d⌋
		dcoef.quo(bpow10[bscale+d.Scale()], dcoef)
		eneg = true
	} else {
		dcoef.lsh(dcoef, bscale-d.Scale())
	}

	// Compute e = log(d)
	ecoef.log(dcoef)

	// Compute e = ⌊e / log(10)⌋
	ecoef.lsh(ecoef, bscale)
	ecoef.quo(ecoef, blog[10])

	return newFromBint(eneg, ecoef, bscale, 0)
}

// Log1p returns the (possibly rounded) shifted natural logarithm of a decimal.
//
// Log1p returns an error if the decimal is equal to or less than negative one.
func (d Decimal) Log1p() (Decimal, error) {
	if d.IsNeg() && d.Cmp(NegOne) <= 0 {
		return Decimal{}, fmt.Errorf("computing log1p(%v): %w: logarithm of a decimal less than or equal to -1", d, errInvalidOperation)
	}

	// Special case: zero
	if d.IsZero() {
		return newSafe(false, 0, 0)
	}

	// General case
	e, err := d.log1pBint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing log1p(%v): %w", d, err)
	}

	return e, nil
}

// log1pBint computes the shifted natural logarithm of a decimal using *big.Int arithmetic.
func (d Decimal) log1pBint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	eneg := false

	// Alignment
	if d.IsNeg() {
		// Compute d = ⌊1 / (d + 1)⌋
		dcoef.subAbs(dcoef, bpow10[d.Scale()])
		dcoef.quo(bpow10[bscale+d.Scale()], dcoef)
		eneg = true
	} else {
		// Compute d = d + 1
		dcoef.add(dcoef, bpow10[d.Scale()])
		dcoef.lsh(dcoef, bscale-d.Scale())
	}

	// Compute e = log(d)
	ecoef.log(dcoef)

	return newFromBint(eneg, ecoef, bscale, 0)
}

// Log returns the (possibly rounded) natural logarithm of a decimal.
//
// Log returns an error if the decimal is zero or negative.
func (d Decimal) Log() (Decimal, error) {
	// Special case: zero or negative
	if !d.IsPos() {
		return Decimal{}, fmt.Errorf("computing log(%v): %w: logarithm of non-positive", d, errInvalidOperation)
	}

	// Special case: one
	if d.IsOne() {
		return newSafe(false, 0, 0)
	}

	// General case
	e, err := d.logBint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing log(%v): %w", d, err)
	}

	return e, nil
}

// logBint computes the natural logarithm of a decimal using *big.Int arithmetic.
func (d Decimal) logBint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	eneg := false

	// Alignment
	if d.WithinOne() {
		// Compute d = ⌊1 / d⌋
		dcoef.quo(bpow10[bscale+d.Scale()], dcoef)
		eneg = true
	} else {
		dcoef.lsh(dcoef, bscale-d.Scale())
	}

	// Compute e = log(d)
	ecoef.log(dcoef)

	return newFromBint(eneg, ecoef, bscale, 0)
}

// log calculates z = log(x) using Halley's method.
// The argument x must satisfy x >= 1, otherwise the result is undefined.
// x must be represented as a big integer: round(x * 10^41).
// The result z is represented as a big integer: round(z * 10^41).
func (z *bint) log(x *bint) {
	zcoef := getBint()
	defer putBint(zcoef)

	fcoef := getBint()
	defer putBint(fcoef)

	Ecoef := getBint()
	defer putBint(Ecoef)

	ncoef := getBint()
	defer putBint(ncoef)

	mcoef := getBint()
	defer putBint(mcoef)

	fcoef.setFint(0)

	// The initial guess is calculated as n*ln(10),
	// where n is the position of the most significant digit.
	n := x.prec() - bscale
	zcoef.setBint(bnlog10[n])

	// Halley's method
	for range 50 {
		Ecoef.exp(zcoef)
		ncoef.sub(Ecoef, x)
		ncoef.dbl(ncoef)
		mcoef.add(Ecoef, x)
		ncoef.lsh(ncoef, bscale)
		ncoef.quo(ncoef, mcoef)
		fcoef.sub(zcoef, ncoef)
		if zcoef.cmp(fcoef) == 0 {
			break
		}
		zcoef.setBint(fcoef)
	}

	z.setBint(zcoef)
}

// Exp returns the (possibly rounded) exponential of a decimal.
//
// Exp returns an error if the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) Exp() (Decimal, error) {
	// Special case: zero
	if d.IsZero() {
		return newSafe(false, 1, 0)
	}

	// Special case: overflow
	if d.CmpAbs(Hundred) >= 0 {
		if !d.IsNeg() {
			return Decimal{}, fmt.Errorf("computing exp(%v): %w", d, unknownOverflowError())
		}
		return newSafe(false, 0, MaxScale)
	}

	// General case
	e, err := d.expBint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing exp(%v): %w", d, err)
	}

	return e, nil
}

// expBint computes exponential of a decimal using *big.Int arithmetic.
func (d Decimal) expBint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)

	// Alignment
	dcoef.lsh(dcoef, bscale-d.Scale())

	// Compute e = exp(d)
	ecoef.exp(dcoef)

	if d.IsNeg() {
		if ecoef.sign() == 0 {
			return Decimal{}, unknownOverflowError()
		}
		// Compute e = ⌊1 / e⌋
		ecoef.quo(bpow10[2*bscale], ecoef)
	}

	return newFromBint(false, ecoef, bscale, 0)
}

// Expm1 returns the (possibly rounded) shifted exponential of a decimal.
//
// Expm1 returns an error if the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) Expm1() (Decimal, error) {
	// Special case: zero
	if d.IsZero() {
		return newSafe(false, 0, 0)
	}

	// Special case: overflow
	if d.CmpAbs(Hundred) >= 0 {
		if !d.IsNeg() {
			return Decimal{}, fmt.Errorf("computing expm1(%v): %w", d, unknownOverflowError())
		}
		return newSafe(true, pow10[MaxScale-1], MaxScale-1)
	}

	// General case
	e, err := d.expm1Bint()
	if err != nil {
		return Decimal{}, fmt.Errorf("computing expm1(%v): %w", d, err)
	}

	return e, nil
}

// expm1Bint computes shifted exponential of a decimal using *big.Int arithmetic.
func (d Decimal) expm1Bint() (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)

	// Alignment
	dcoef.lsh(dcoef, bscale-d.Scale())

	// Compute e = exp(d)
	ecoef.exp(dcoef)

	if d.IsNeg() {
		if ecoef.sign() == 0 {
			return Decimal{}, unknownOverflowError()
		}
		// Compute e = ⌊1 / e⌋
		ecoef.quo(bpow10[2*bscale], ecoef)
	}

	// Compute e = e - 1
	eneg := false
	if ecoef.cmp(bpow10[bscale]) < 0 {
		eneg = true
	}
	ecoef.subAbs(ecoef, bpow10[bscale])

	return newFromBint(eneg, ecoef, bscale, 0)
}

// exp calculates z = exp(x) using Taylor series expansion.
// The argument x must satisfy 0 <= x < 100, otherwise the result is undefined.
// The argument x must be represented as a big integer: round(x * 10^41).
// The result z is represented as a big integer: round(z * 10^41).
func (z *bint) exp(x *bint) {
	qcoef := getBint()
	defer putBint(qcoef)

	rcoef := getBint()
	defer putBint(rcoef)

	// Split x into integer part q and fractional part r
	qcoef.quoRem(x, bpow10[bscale], rcoef)

	// Retrieve z = exp(q) from precomputed cache
	z.setBint(bexp[int(qcoef.fint())]) //nolint:gosec

	if rcoef.sign() == 0 {
		return
	}

	zcoef := getBint()
	defer putBint(zcoef)

	gcoef := getBint()
	defer putBint(gcoef)

	hcoef := getBint()
	defer putBint(hcoef)

	zcoef.setFint(0)
	gcoef.setBint(bpow10[bscale])

	// Compute exp(r) using Taylor series expansion
	// exp(r) = r^0 / 0! + r^1 / 1! + ... + r^n / n!
	for i := range len(bfact) {
		hcoef.quo(gcoef, bfact[i])
		if hcoef.sign() == 0 {
			break
		}
		zcoef.add(zcoef, hcoef)
		gcoef.mul(gcoef, rcoef)
		gcoef.rshDown(gcoef, bscale)
	}

	// Compute z = z * exp(r)
	z.mul(z, zcoef)
	z.rshDown(z, bscale)
}

// Sum returns the (possibly rounded) sum of decimals.
// It computes d1 + d2 + ... + dn without intermediate rounding.
//
// Sum returns an error if:
//   - no argements are provided;
//   - the integer part of the result has more than [MaxPrec] digits.
func Sum(d ...Decimal) (Decimal, error) {
	// Special cases
	switch len(d) {
	case 0:
		return Decimal{}, fmt.Errorf("computing [sum([])]: %w", errInvalidOperation)
	case 1:
		return d[0], nil
	}

	// General case
	e, err := sumFint(d...)
	if err != nil {
		e, err = sumBint(d...)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [sum(%v)]: %w", d, err)
		}
	}

	return e, nil
}

// sumFint computes the sum of decimals using uint64 arithmetic.
func sumFint(d ...Decimal) (Decimal, error) {
	ecoef := Zero.coef
	escale := Zero.Scale()
	eneg := Zero.IsNeg()

	for _, f := range d {
		fcoef := f.coef

		// Alignment
		var ok bool
		switch {
		case escale > f.Scale():
			fcoef, ok = fcoef.lsh(escale - f.Scale())
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
		case escale < f.Scale():
			ecoef, ok = ecoef.lsh(f.Scale() - escale)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
			escale = f.Scale()
		}

		// Compute e = e + f
		if eneg == f.IsNeg() {
			ecoef, ok = ecoef.add(fcoef)
			if !ok {
				return Decimal{}, errDecimalOverflow
			}
		} else {
			if fcoef > ecoef {
				eneg = f.IsNeg()
			}
			ecoef = ecoef.subAbs(fcoef)
		}
	}

	return newFromFint(eneg, ecoef, escale, 0)
}

// sumBint computes the sum of decimals using *big.Int arithmetic.
func sumBint(d ...Decimal) (Decimal, error) {
	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	ecoef.setFint(Zero.coef)
	escale := Zero.Scale()
	eneg := Zero.IsNeg()

	for _, f := range d {
		fcoef.setFint(f.coef)

		// Alignment
		switch {
		case escale > f.Scale():
			fcoef.lsh(fcoef, escale-f.Scale())
		case escale < f.Scale():
			ecoef.lsh(ecoef, f.Scale()-escale)
			escale = f.Scale()
		}

		// Compute e = e + f
		if eneg == f.IsNeg() {
			ecoef.add(ecoef, fcoef)
		} else {
			if fcoef.cmp(ecoef) > 0 {
				eneg = f.IsNeg()
			}
			ecoef.subAbs(ecoef, fcoef)
		}
	}

	return newFromBint(eneg, ecoef, escale, 0)
}

// SubAbs returns the (possibly rounded) absolute difference between decimals d and e.
//
// SubAbs returns an error if the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) SubAbs(e Decimal) (Decimal, error) {
	f, err := d.Sub(e)
	if err != nil {
		return Decimal{}, fmt.Errorf("computing [abs(%v - %v)]: %w", d, e, err)
	}
	return f.Abs(), nil
}

// Sub returns the (possibly rounded) difference between decimals d and e.
//
// Sub returns an error if the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) Sub(e Decimal) (Decimal, error) {
	return d.AddExact(e.Neg(), 0)
}

// SubExact is similar to [Decimal.Sub], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) SubExact(e Decimal, scale int) (Decimal, error) {
	return d.AddExact(e.Neg(), scale)
}

// Add returns the (possibly rounded) sum of decimals d and e.
//
// Add returns an error if the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) Add(e Decimal) (Decimal, error) {
	return d.AddExact(e, 0)
}

// AddExact is similar to [Decimal.Add], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) AddExact(e Decimal, scale int) (Decimal, error) {
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("computing [%v + %v]: %w", d, e, errScaleRange)
	}

	// General case
	f, err := d.addFint(e, scale)
	if err != nil {
		f, err = d.addBint(e, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v + %v]: %w", d, e, err)
		}
	}

	return f, nil
}

// addFint computes the sum of two decimals using uint64 arithmetic.
func (d Decimal) addFint(e Decimal, minScale int) (Decimal, error) {
	dcoef := d.coef
	dscale := d.Scale()
	dneg := d.IsNeg()

	ecoef := e.coef

	// Alignment
	var ok bool
	switch {
	case dscale > e.Scale():
		ecoef, ok = ecoef.lsh(dscale - e.Scale())
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	case dscale < e.Scale():
		dcoef, ok = dcoef.lsh(e.Scale() - dscale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		dscale = e.Scale()
	}

	// Compute d = d + e
	if dneg == e.IsNeg() {
		dcoef, ok = dcoef.add(ecoef)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	} else {
		if ecoef > dcoef {
			dneg = e.IsNeg()
		}
		dcoef = dcoef.subAbs(ecoef)
	}

	return newFromFint(dneg, dcoef, dscale, minScale)
}

// addBint computes the sum of two decimals using *big.Int arithmetic.
func (d Decimal) addBint(e Decimal, minScale int) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	dscale := d.Scale()
	ecoef.setFint(e.coef)
	dneg := d.IsNeg()

	// Alignment
	switch {
	case dscale > e.Scale():
		ecoef.lsh(ecoef, dscale-e.Scale())
	case dscale < e.Scale():
		dcoef.lsh(dcoef, e.Scale()-dscale)
		dscale = e.Scale()
	}

	// Compute d = d + e
	if dneg == e.IsNeg() {
		dcoef.add(dcoef, ecoef)
	} else {
		if ecoef.cmp(dcoef) > 0 {
			dneg = e.IsNeg()
		}
		dcoef.subAbs(dcoef, ecoef)
	}

	return newFromBint(dneg, dcoef, dscale, minScale)
}

// SubMul returns the (possibly rounded) [fused multiply-subtraction] of decimals d, e, and f.
// It computes d - e * f without any intermediate rounding.
// This method is useful for improving the accuracy and performance of algorithms
// that involve the accumulation of products, such as daily interest accrual.
//
// SubMul returns an error if the integer part of the result has more than [MaxPrec] digits.
//
// [fused multiply-subtraction]: https://en.wikipedia.org/wiki/Multiply%E2%80%93accumulate_operation#Fused_multiply%E2%80%93add
func (d Decimal) SubMul(e, f Decimal) (Decimal, error) {
	return d.AddMulExact(e.Neg(), f, 0)
}

// SubMulExact is similar to [Decimal.SubMul], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) SubMulExact(e, f Decimal, scale int) (Decimal, error) {
	return d.AddMulExact(e.Neg(), f, scale)
}

// AddMul returns the (possibly rounded) [fused multiply-addition] of decimals d, e, and f.
// It computes d + e * f without any intermediate rounding.
// This method is useful for improving the accuracy and performance of algorithms
// that involve the accumulation of products, such as daily interest accrual.
//
// AddMul returns an error if the integer part of the result has more than [MaxPrec] digits.
//
// [fused multiply-addition]: https://en.wikipedia.org/wiki/Multiply%E2%80%93accumulate_operation#Fused_multiply%E2%80%93add
func (d Decimal) AddMul(e, f Decimal) (Decimal, error) {
	return d.AddMulExact(e, f, 0)
}

// AddMulExact is similar to [Decimal.AddMul], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) AddMulExact(e, f Decimal, scale int) (Decimal, error) {
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("computing [%v + %v * %v]: %w", d, e, f, errScaleRange)
	}

	// General case
	g, err := d.addMulFint(e, f, scale)
	if err != nil {
		g, err = d.addMulBint(e, f, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v + %v * %v]: %w", d, e, f, err)
		}
	}

	return g, nil
}

// addMulFint computes the fused multiply-addition of three decimals using uint64 arithmetic.
func (d Decimal) addMulFint(e, f Decimal, minScale int) (Decimal, error) {
	dcoef := d.coef
	dscale := d.Scale()
	dneg := d.IsNeg()

	ecoef := e.coef
	escale := e.Scale()
	eneg := e.IsNeg()

	fcoef := f.coef

	// Compute e = e * f
	var ok bool
	ecoef, ok = ecoef.mul(fcoef)
	if !ok {
		return Decimal{}, errDecimalOverflow
	}
	escale = escale + f.Scale()
	eneg = eneg != f.IsNeg()

	// Alignment
	switch {
	case dscale > escale:
		ecoef, ok = ecoef.lsh(dscale - escale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	case dscale < escale:
		dcoef, ok = dcoef.lsh(escale - dscale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		dscale = escale
	}

	// Compute d = d + e
	if dneg == eneg {
		dcoef, ok = dcoef.add(ecoef)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	} else {
		if ecoef > dcoef {
			dneg = eneg
		}
		dcoef = dcoef.subAbs(ecoef)
	}

	return newFromFint(dneg, dcoef, dscale, minScale)
}

// addMulBint computes the fused multiply-addition of three decimals using *big.Int arithmetic.
func (d Decimal) addMulBint(e, f Decimal, minScale int) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	dcoef.setFint(d.coef)
	dscale := d.Scale()
	dneg := d.IsNeg()
	ecoef.setFint(e.coef)
	escale := e.Scale()
	eneg := e.IsNeg()
	fcoef.setFint(f.coef)

	// Compute e = e * f
	ecoef.mul(ecoef, fcoef)
	escale = escale + f.Scale()
	eneg = eneg != f.IsNeg()

	// Alignment
	switch {
	case dscale > escale:
		ecoef.lsh(ecoef, dscale-escale)
	case dscale < escale:
		dcoef.lsh(dcoef, escale-d.Scale())
		dscale = escale
	}

	// Compute d = d + e
	if dneg == eneg {
		dcoef.add(dcoef, ecoef)
	} else {
		if ecoef.cmp(dcoef) > 0 {
			dneg = eneg
		}
		dcoef.subAbs(dcoef, ecoef)
	}

	return newFromBint(dneg, dcoef, dscale, minScale)
}

// SubQuo returns the (possibly rounded) fused quotient-subtraction of decimals d, e, and f.
// It computes d - e / f with at least double precision during intermediate rounding.
// This method is useful for improving the accuracy and performance of algorithms
// that involve the accumulation of quotients, such as internal rate of return.
//
// AddQuo returns an error if:
//   - the divisor is 0;
//   - the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) SubQuo(e, f Decimal) (Decimal, error) {
	return d.AddQuoExact(e.Neg(), f, 0)
}

// SubQuoExact is similar to [Decimal.SubQuo], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) SubQuoExact(e, f Decimal, scale int) (Decimal, error) {
	return d.AddQuoExact(e.Neg(), f, scale)
}

// AddQuo returns the (possibly rounded) fused quotient-addition of decimals d, e, and f.
// It computes d + e / f with at least double precision during the intermediate rounding.
// This method is useful for improving the accuracy and performance of algorithms
// that involve the accumulation of quotients, such as internal rate of return.
//
// AddQuo returns an error if:
//   - the divisor is 0;
//   - the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) AddQuo(e, f Decimal) (Decimal, error) {
	return d.AddQuoExact(e, f, 0)
}

// AddQuoExact is similar to [Decimal.AddQuo], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) AddQuoExact(e, f Decimal, scale int) (Decimal, error) {
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("computing [%v + %v / %v]: %w", d, e, f, errScaleRange)
	}

	// Special case: zero divisor
	if f.IsZero() {
		return Decimal{}, fmt.Errorf("computing [%v + %v / %v]: %w", d, e, f, errDivisionByZero)
	}

	// Special case: zero dividend
	if e.IsZero() {
		scale = max(scale, e.Scale()-f.Scale())
		return d.Pad(scale), nil
	}

	// General case
	g, err := d.addQuoFint(e, f, scale)
	if err != nil {
		g, err = d.addQuoBint(e, f, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v + %v / %v]: %w", d, e, f, err)
		}
	}

	// Preferred scale
	scale = max(scale, d.Scale(), e.Scale()-f.Scale())
	g = g.Trim(scale)

	return g, nil
}

// addQuoFint computes the fused quotient-addition of three decimals using uint64 arithmetic.
func (d Decimal) addQuoFint(e, f Decimal, minScale int) (Decimal, error) {
	dcoef := d.coef
	dscale := d.Scale()
	dneg := d.IsNeg()

	ecoef := e.coef
	escale := e.Scale()
	eneg := e.IsNeg()

	fcoef := f.coef

	// Alignment
	var ok bool
	if shift := MaxPrec - ecoef.prec(); shift > 0 {
		ecoef, ok = ecoef.lsh(shift)
		if !ok {
			return Decimal{}, errDecimalOverflow // Should never happen
		}
		escale = escale + shift
	}
	if shift := fcoef.ntz(); shift > 0 {
		fcoef = fcoef.rshDown(shift)
		escale = escale + shift
	}

	// Compute e = e / f
	ecoef, ok = ecoef.quo(fcoef)
	if !ok {
		return Decimal{}, errInexactDivision
	}
	escale = escale - f.Scale()
	eneg = eneg != f.IsNeg()

	// Alignment
	switch {
	case dscale > escale:
		ecoef, ok = ecoef.lsh(dscale - escale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	case dscale < escale:
		if shift := min(escale-e.Scale()+f.Scale(), escale-dscale, ecoef.ntz()); shift > 0 {
			ecoef = ecoef.rshDown(shift)
			escale = escale - shift
		}
		dcoef, ok = dcoef.lsh(escale - dscale)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
		dscale = escale
	}

	// Compute d = d + e
	if dneg == eneg {
		dcoef, ok = dcoef.add(ecoef)
		if !ok {
			return Decimal{}, errDecimalOverflow
		}
	} else {
		if ecoef > dcoef {
			dneg = eneg
		}
		dcoef = dcoef.subAbs(ecoef)
	}

	return newFromFint(dneg, dcoef, dscale, minScale)
}

// addQuoBint computes the fused quotient-addition of three decimals using *big.Int arithmetic.
func (d Decimal) addQuoBint(e, f Decimal, minScale int) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	fcoef := getBint()
	defer putBint(fcoef)

	dcoef.setFint(d.coef)
	dneg := d.IsNeg()
	ecoef.setFint(e.coef)
	eneg := e.IsNeg()
	fcoef.setFint(f.coef)

	// Alignment
	ecoef.lsh(ecoef, bscale-e.Scale()+f.Scale())

	// Compute e = ⌊e / f⌋
	ecoef.quo(ecoef, fcoef)
	eneg = eneg != f.IsNeg()

	// Alignment
	dcoef.lsh(dcoef, bscale-d.Scale())

	// Compute d = d + e
	if dneg == eneg {
		dcoef.add(dcoef, ecoef)
	} else {
		if ecoef.cmp(dcoef) > 0 {
			dneg = eneg
		}
		dcoef.subAbs(dcoef, ecoef)
	}

	return newFromBint(dneg, dcoef, bscale, minScale)
}

// Inv returns the (possibly rounded) inverse of the decimal.
//
// Inv returns an error if:
//   - the integer part of the result has more than [MaxPrec] digits;
//   - the decimal is 0.
func (d Decimal) Inv() (Decimal, error) {
	f, err := One.Quo(d)
	if err != nil {
		return Decimal{}, fmt.Errorf("inverting %v: %w", d, err)
	}
	return f, nil
}

// Quo returns the (possibly rounded) quotient of decimals d and e.
//
// Quo returns an error if:
//   - the divisor is 0;
//   - the integer part of the result has more than [MaxPrec] digits.
func (d Decimal) Quo(e Decimal) (Decimal, error) {
	return d.QuoExact(e, 0)
}

// QuoExact is similar to [Decimal.Quo], but it allows you to specify the number of digits
// after the decimal point that should be considered significant.
// If any of the significant digits are lost during rounding, the method will return an error.
// This method is useful for financial calculations where the scale should be
// equal to or greater than the currency's scale.
func (d Decimal) QuoExact(e Decimal, scale int) (Decimal, error) {
	if scale < MinScale || scale > MaxScale {
		return Decimal{}, fmt.Errorf("computing [%v / %v]: %w", d, e, errScaleRange)
	}

	// Special case: zero divisor
	if e.IsZero() {
		return Decimal{}, fmt.Errorf("computing [%v / %v]: %w", d, e, errDivisionByZero)
	}

	// Special case: zero dividend
	if d.IsZero() {
		scale = max(scale, d.Scale()-e.Scale())
		return newSafe(false, 0, scale)
	}

	// General case
	f, err := d.quoFint(e, scale)
	if err != nil {
		f, err = d.quoBint(e, scale)
		if err != nil {
			return Decimal{}, fmt.Errorf("computing [%v / %v]: %w", d, e, err)
		}
	}

	// Preferred scale
	scale = max(scale, d.Scale()-e.Scale())
	f = f.Trim(scale)

	return f, nil
}

// quoFint computes the quotient of two decimals using uint64 arithmetic.
func (d Decimal) quoFint(e Decimal, minScale int) (Decimal, error) {
	dcoef := d.coef
	dscale := d.Scale()
	dneg := d.IsNeg()

	ecoef := e.coef

	// Alignment
	var ok bool
	if shift := MaxPrec - dcoef.prec(); shift > 0 {
		dcoef, ok = dcoef.lsh(shift)
		if !ok {
			return Decimal{}, errDecimalOverflow // Should never happen
		}
		dscale = dscale + shift
	}
	if shift := ecoef.ntz(); shift > 0 {
		ecoef = ecoef.rshDown(shift)
		dscale = dscale + shift
	}

	// Compute d = d / e
	dcoef, ok = dcoef.quo(ecoef)
	if !ok {
		return Decimal{}, errInexactDivision
	}
	dscale = dscale - e.Scale()
	dneg = dneg != e.IsNeg()

	return newFromFint(dneg, dcoef, dscale, minScale)
}

// quoBint computes the quotient of two decimals using *big.Int arithmetic.
func (d Decimal) quoBint(e Decimal, minScale int) (Decimal, error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	dneg := d.IsNeg()
	ecoef.setFint(e.coef)

	// Alignment
	dcoef.lsh(dcoef, bscale+e.Scale()-d.Scale())

	// Compute d = ⌊d / e⌋
	dcoef.quo(dcoef, ecoef)
	dneg = dneg != e.IsNeg()

	return newFromBint(dneg, dcoef, bscale, minScale)
}

// QuoRem returns the quotient q and remainder r of decimals d and e
// such that d = e * q + r, where q is an integer and the sign of the
// reminder r is the same as the sign of the dividend d.
//
// QuoRem returns an error if:
//   - the divisor is 0;
//   - the integer part of the quotient has more than [MaxPrec] digits.
func (d Decimal) QuoRem(e Decimal) (q, r Decimal, err error) {
	// Special case: zero divisor
	if e.IsZero() {
		return Decimal{}, Decimal{}, fmt.Errorf("computing [%v div %v] and [%v mod %v]: %w", d, e, d, e, errDivisionByZero)
	}

	// General case
	q, r, err = d.quoRemFint(e)
	if err != nil {
		q, r, err = d.quoRemBint(e)
		if err != nil {
			return Decimal{}, Decimal{}, fmt.Errorf("computing [%v div %v] and [%v mod %v]: %w", d, e, d, e, err)
		}
	}

	return q, r, nil
}

// quoRemFint computes the quotient and remainder of two decimals using uint64 arithmetic.
func (d Decimal) quoRemFint(e Decimal) (q, r Decimal, err error) {
	dcoef := d.coef
	ecoef := e.coef
	rscale := d.Scale()

	// Alignment
	var ok bool
	switch {
	case d.Scale() > e.Scale():
		ecoef, ok = ecoef.lsh(d.Scale() - e.Scale())
		if !ok {
			return Decimal{}, Decimal{}, errDecimalOverflow
		}
	case d.Scale() < e.Scale():
		dcoef, ok = dcoef.lsh(e.Scale() - d.Scale())
		if !ok {
			return Decimal{}, Decimal{}, errDecimalOverflow
		}
		rscale = e.Scale()
	}

	// Compute q = ⌊d / e⌋, r = d - e * q
	qcoef, rcoef, ok := dcoef.quoRem(ecoef)
	if !ok {
		return Decimal{}, Decimal{}, errDivisionByZero // Should never happen
	}
	qsign := d.IsNeg() != e.IsNeg()
	rsign := d.IsNeg()

	q, err = newFromFint(qsign, qcoef, 0, 0)
	if err != nil {
		return Decimal{}, Decimal{}, err
	}
	r, err = newFromFint(rsign, rcoef, rscale, rscale)
	if err != nil {
		return Decimal{}, Decimal{}, err
	}
	return q, r, nil
}

// quoRemBint computes the quotient and remainder of two decimals using *big.Int arithmetic.
func (d Decimal) quoRemBint(e Decimal) (q, r Decimal, err error) {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	qcoef := getBint()
	defer putBint(qcoef)

	rcoef := getBint()
	defer putBint(rcoef)

	dcoef.setFint(d.coef)
	ecoef.setFint(e.coef)
	rscale := d.Scale()

	// Alignment
	switch {
	case d.Scale() > e.Scale():
		ecoef.lsh(ecoef, d.Scale()-e.Scale())
	case d.Scale() < e.Scale():
		dcoef.lsh(dcoef, e.Scale()-d.Scale())
		rscale = e.Scale()
	}

	// Compute q = ⌊d / e⌋, r = d - e * q
	qcoef.quoRem(dcoef, ecoef, rcoef)
	qsign := d.IsNeg() != e.IsNeg()
	rsign := d.IsNeg()

	q, err = newFromBint(qsign, qcoef, 0, 0)
	if err != nil {
		return Decimal{}, Decimal{}, err
	}
	r, err = newFromBint(rsign, rcoef, rscale, rscale)
	if err != nil {
		return Decimal{}, Decimal{}, err
	}
	return q, r, nil
}

// Max returns the larger decimal.
// See also method [Decimal.CmpTotal].
func (d Decimal) Max(e Decimal) Decimal {
	if d.CmpTotal(e) >= 0 {
		return d
	}
	return e
}

// Min returns the smaller decimal.
// See also method [Decimal.CmpTotal].
func (d Decimal) Min(e Decimal) Decimal {
	if d.CmpTotal(e) <= 0 {
		return d
	}
	return e
}

// Clamp compares decimals and returns:
//
//	min if d < min
//	max if d > max
//	  d otherwise
//
// See also method [Decimal.CmpTotal].
//
// Clamp returns an error if min is greater than max numerically.
//
//nolint:revive
func (d Decimal) Clamp(min, max Decimal) (Decimal, error) {
	if min.Cmp(max) > 0 {
		return Decimal{}, fmt.Errorf("clamping %v: invalid range", d)
	}
	if min.CmpTotal(max) > 0 {
		// min and max are equal numerically but have different scales.
		// Swaping min and max to ensure total ordering.
		min, max = max, min
	}
	if d.CmpTotal(min) < 0 {
		return min, nil
	}
	if d.CmpTotal(max) > 0 {
		return max, nil
	}
	return d, nil
}

// CmpTotal compares decimal representations and returns:
//
//	-1 if d < e
//	-1 if d = e and d.scale > e.scale
//	 0 if d = e and d.scale = e.scale
//	+1 if d = e and d.scale < e.scale
//	+1 if d > e
//
// See also method [Decimal.Cmp].
func (d Decimal) CmpTotal(e Decimal) int {
	switch d.Cmp(e) {
	case -1:
		return -1
	case 1:
		return 1
	}
	switch {
	case d.Scale() > e.Scale():
		return -1
	case d.Scale() < e.Scale():
		return 1
	}
	return 0
}

// CmpAbs compares absolute values of decimals and returns:
//
//	-1 if |d| < |e|
//	 0 if |d| = |e|
//	+1 if |d| > |e|
//
// See also method [Decimal.Cmp].
func (d Decimal) CmpAbs(e Decimal) int {
	d, e = d.Abs(), e.Abs()
	return d.Cmp(e)
}

// Equal compares decimals and returns:
//
//	 true if d = e
//	false otherwise
//
// See also method [Decimal.Cmp].
func (d Decimal) Equal(e Decimal) bool {
	return d.Cmp(e) == 0
}

// Less compares decimals and returns:
//
//	 true if d < e
//	false otherwise
//
// See also method [Decimal.Cmp].
func (d Decimal) Less(e Decimal) bool {
	return d.Cmp(e) < 0
}

// Cmp compares decimals and returns:
//
//	-1 if d < e
//	 0 if d = e
//	+1 if d > e
//
// See also methods [Decimal.CmpAbs], [Decimal.CmpTotal].
func (d Decimal) Cmp(e Decimal) int {
	// Special case: different signs
	switch {
	case d.Sign() > e.Sign():
		return 1
	case d.Sign() < e.Sign():
		return -1
	}

	// General case
	r, err := d.cmpFint(e)
	if err != nil {
		r = d.cmpBint(e)
	}
	return r
}

// cmpFint compares decimals using uint64 arithmetic.
func (d Decimal) cmpFint(e Decimal) (int, error) {
	dcoef := d.coef
	ecoef := e.coef

	// Alignment
	var ok bool
	switch {
	case d.Scale() > e.Scale():
		ecoef, ok = ecoef.lsh(d.Scale() - e.Scale())
		if !ok {
			return 0, errDecimalOverflow
		}
	case d.Scale() < e.Scale():
		dcoef, ok = dcoef.lsh(e.Scale() - d.Scale())
		if !ok {
			return 0, errDecimalOverflow
		}
	}

	// Comparison
	switch {
	case dcoef > ecoef:
		return d.Sign(), nil
	case ecoef > dcoef:
		return -e.Sign(), nil
	}
	return 0, nil
}

// cmpBint compares decimals using *big.Int arithmetic.
func (d Decimal) cmpBint(e Decimal) int {
	dcoef := getBint()
	defer putBint(dcoef)

	ecoef := getBint()
	defer putBint(ecoef)

	dcoef.setFint(d.coef)
	ecoef.setFint(e.coef)

	// Alignment
	switch {
	case d.Scale() > e.Scale():
		ecoef.lsh(ecoef, d.Scale()-e.Scale())
	case d.Scale() < e.Scale():
		dcoef.lsh(dcoef, e.Scale()-d.Scale())
	}

	// Comparison
	switch dcoef.cmp(ecoef) {
	case 1:
		return d.Sign()
	case -1:
		return -e.Sign()
	}
	return 0
}

// NullDecimal represents a decimal that can be null.
// Its zero value is null.
// NullDecimal is not thread-safe.
type NullDecimal struct {
	Decimal Decimal
	Valid   bool
}

// Scan implements the [sql.Scanner] interface.
// See also method [Decimal.Scan].
//
// [sql.Scanner]: https://pkg.go.dev/database/sql#Scanner
func (n *NullDecimal) Scan(value any) error {
	if value == nil {
		n.Decimal = Decimal{}
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.Decimal.Scan(value)
}

// Value implements the [driver.Valuer] interface.
// See also method [Decimal.Value].
//
// [driver.Valuer]: https://pkg.go.dev/database/sql/driver#Valuer
func (n NullDecimal) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Decimal.Value()
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
// See also method [Decimal.UnmarshalJSON].
//
// [json.Unmarshaler]: https://pkg.go.dev/encoding/json#Unmarshaler
func (n *NullDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Decimal = Decimal{}
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.Decimal.UnmarshalJSON(data)
}

// MarshalJSON implements the [json.Marshaler] interface.
// See also method [Decimal.MarshalJSON].
//
// [json.Marshaler]: https://pkg.go.dev/encoding/json#Marshaler
func (n NullDecimal) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return n.Decimal.MarshalJSON()
}

// UnmarshalBSONValue implements the [v2/bson.ValueUnmarshaler] interface.
// UnmarshalBSONValue supports the following [types]: Null, Double, String, 32-bit Integer, 64-bit Integer, and [Decimal128].
// See also method [Decimal.UnmarshalBSONValue].
//
// [v2/bson.ValueUnmarshaler]: https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/bson#ValueUnmarshaler
// [types]: https://bsonspec.org/spec.html
// [Decimal128]: https://github.com/mongodb/specifications/blob/master/source/bson-decimal128/decimal128.md
func (n *NullDecimal) UnmarshalBSONValue(typ byte, data []byte) error {
	// constants are from https://bsonspec.org/spec.html
	if typ == 10 {
		n.Decimal = Decimal{}
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.Decimal.UnmarshalBSONValue(typ, data)
}

// MarshalBSONValue implements the [v2/bson.ValueMarshaler] interface.
// MarshalBSONValue returns [Null] or [Decimal128].
// See also method [Decimal.MarshalBSONValue].
//
// [v2/bson.ValueMarshaler]: https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/bson#ValueMarshaler
// [Null]: https://bsonspec.org/spec.html
// [Decimal128]: https://github.com/mongodb/specifications/blob/master/source/bson-decimal128/decimal128.md
func (n NullDecimal) MarshalBSONValue() (typ byte, data []byte, err error) {
	// constants are from https://bsonspec.org/spec.html
	if !n.Valid {
		return 10, nil, nil
	}
	return n.Decimal.MarshalBSONValue()
}
