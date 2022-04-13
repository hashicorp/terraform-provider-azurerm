package datetime

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

const nsecsPerSec = 1000000000

var zeroTime = time.Time{}

type parser struct {
	s   *scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

func newParser(r io.Reader) *parser {
	return &parser{s: newScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *parser) scan() (tok token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *parser) unscan() { p.buf.n = 1 }

func (p *parser) parse(defaultLocation *time.Location) (time.Time, error) {
	location := defaultLocation
	year, month, day, err := p.parseDate()
	if err != nil {
		return zeroTime, err
	}

	var hour, min, sec, nsec int
	// parseDate has already checked that there's a T or EOF next, so those are the only cases we need
	// to check here.

	switch tok, _ := p.scan(); tok {
	case EOF:
		return buildTime(year, month, day, hour, min, sec, nsec, location)
	case T:
		hour, min, sec, nsec, err = p.parseTime()
		if err != nil {
			return zeroTime, err
		}
	}

	location, err = p.parseLocation(location)
	if err != nil {
		return zeroTime, err
	}

	// there should be nothing left at this point
	if tok, lit := p.scan(); tok != EOF {
		return zeroTime, fmt.Errorf("expected EOF. got %s", lit)
	}

	return buildTime(year, month, day, hour, min, sec, nsec, location)
}

func (p *parser) parseLocation(defaultLocation *time.Location) (*time.Location, error) {
	var sign, secs int
	var name string
	switch tok, lit := p.scan(); tok {
	case EOF:
		return defaultLocation, nil
	case Z:
		return time.UTC, nil
	case PLUS:
		sign = 1
		name += lit
	case DASH:
		sign = -1
		name += lit
	default:
		return nil, fmt.Errorf("expected Z, timezone offset, or EOF. got %s", lit)
	}

	switch tok, lit := p.scan(); tok {
	case NUMBER:
		name += lit
		switch len(lit) {
		case 4:
			hours := parseInt(lit[:2])
			secs = hours * 60 * 60
			minutes := parseInt(lit[2:4])
			secs += minutes * 60
		case 2:
			hours := parseInt(lit)
			secs = hours * 60 * 60
		default:
			return nil, fmt.Errorf("expected ±hh:mm, ±hhmm, or ±hh timezone offset format. got %s", lit)
		}
	default:
		return nil, fmt.Errorf("expected number. got %s", lit)
	}

	switch tok, lit := p.scan(); tok {
	case EOF:
		return time.FixedZone(name, sign*secs), nil
	default:
		if tok != COLON {
			return nil, fmt.Errorf("expected colon or EOF. got %s", lit)
		}
	}
	name += ":"

	switch tok, lit := p.scan(); tok {
	case NUMBER:
		name += lit
		minutes := parseInt(lit)
		secs += minutes * 60
	default:
		return nil, fmt.Errorf("expected number. got %s", lit)
	}

	return time.FixedZone(name, sign*secs), nil
}

// parseTime returns the hour, minute, second, and nanosecond in the timestamp.
func (p *parser) parseTime() (int, int, int, int, error) {
	var minParsed, secParsed bool
	var hour, min, sec, nsec int
	var err error
	parseErr := func(err error) (int, int, int, int, error) {
		return 0, 0, 0, 0, err
	}

	switch tok, lit := p.scan(); tok {
	case NUMBER:
		// could be hh, hhmm, or hhmmss
		switch len(lit) {
		case 2:
			hour = parseInt(lit)
		case 4:
			hour = parseInt(lit[:2])
			min = parseInt(lit[2:4])
			minParsed = true
		case 6:
			hour = parseInt(lit[:2])
			min = parseInt(lit[2:4])
			minParsed = true
			sec = parseInt(lit[4:6])
			secParsed = true
		default:
			return parseErr(fmt.Errorf("expected time. got %s", lit))
		}
	default:
		return parseErr(fmt.Errorf("expected number. got %s", lit))
	}

	// get the min
	if !minParsed {
		switch tok, lit := p.scan(); tok {
		case EOF:
			return hour, min, sec, nsec, nil
		case COLON:
			min, err = p.scanNumber()
			if err != nil {
				return parseErr(err)
			}
		default:
			if beginsOffset(tok) {
				p.unscan()
				return hour, min, sec, nsec, nil
			}
			return parseErr(fmt.Errorf("expected colon, EOF, or timezone offset. got %s", lit))
		}
	}

	// get the sec
	if !secParsed {
		switch tok, _ := p.scan(); tok {
		case EOF:
			return hour, min, sec, nsec, nil
		case COLON:
			sec, err = p.scanNumber()
			if err != nil {
				return parseErr(err)
			}
		}
	}

	// get the nsec
	switch tok, lit := p.scan(); tok {
	case EOF:
		return hour, min, sec, nsec, nil
	case DOT:
		// can't use scanNumber on the fractional part because we need to preserve leading zeros.
		tok, lit = p.scan()

		if tok != NUMBER {
			return parseErr(fmt.Errorf("expected fraction of seconds.  got %s", lit))
		}

		secFrac := parseDecimal(lit)
		nsec = int(round(secFrac * nsecsPerSec))
	default:
		p.unscan()
	}
	return hour, min, sec, nsec, nil
}

func (p *parser) parseDate() (int, time.Month, int, error) {
	var year int
	month := time.January
	day := 1
	var err error

	// error helpers
	parseErr := func(err error) (int, time.Month, int, error) {
		return 0, time.Month(0), 0, err
	}

	unexpected := func(found, expected string) (int, time.Month, int, error) {
		return parseErr(fmt.Errorf("found %s, expected %s", found, expected))
	}

	// should start with a number like yyyy or yyyymmdd
	switch tok, lit := p.scan(); tok {
	case NUMBER:
		switch len(lit) {
		case 4:
			year = parseInt(lit)
		case 8:
			// we should have yyyymmdd
			year = parseInt(lit[:4])
			monthNum := parseInt(lit[4:6])
			month = time.Month(monthNum)
			day = parseInt(lit[6:8])
			return year, month, day, nil
		default:
			return unexpected(lit, "yyyy-mm-dd or yyyymmdd")
		}
	default:
		return unexpected(lit, "number")
	}

	// if we're here, then we've got a year but not yet a month or day.  Dash or "T" is next.
	switch tok, lit := p.scan(); tok {
	case T, EOF:
		if tok == T {
			p.unscan()
		}
		return year, month, day, nil
	default:
		if tok != DASH {
			return unexpected(lit, "dash, T, or EOF")
		}
	}

	monthNum, err := p.scanNumber()
	if err != nil {
		return parseErr(err)
	}
	month = time.Month(monthNum)

	// if we're here, then we've got a year and month but not yet a day.  Dash or "T" is next.
	switch tok, lit := p.scan(); tok {
	case T, EOF:
		return year, month, day, nil
	default:
		if tok != DASH {
			return unexpected(lit, "dash, T, or EOF")
		}
	}

	day, err = p.scanNumber()
	if err != nil {
		return parseErr(err)
	}

	switch tok, lit := p.scan(); tok {
	case T, EOF:
		if tok == T {
			p.unscan()
		}
	default:
		return unexpected(lit, "T or EOF")
	}

	return year, month, day, nil
}

func (p *parser) scanNumber() (int, error) {
	if tok, lit := p.scan(); tok == NUMBER {
		return strconv.Atoi(lit)
	} else {
		return 0, fmt.Errorf("found %s, expected number", lit)
	}
}

func buildTime(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) (time.Time, error) {
	// year has already been checked for 4 digit length.

	if !checkMonth(int(month)) {
		return zeroTime, fmt.Errorf("%02d is not a valid month", month)
	}

	if !checkDay(month, day) {
		return zeroTime, fmt.Errorf("%02d is not a valid day in %s", day, month)
	}

	if !checkHour(hour) {
		return zeroTime, fmt.Errorf("%02d is not a valid hour", hour)
	}

	if !checkMinSec(min) {
		return zeroTime, fmt.Errorf("%02d is not a valid minute", min)
	}

	if !checkMinSec(sec) {
		return zeroTime, fmt.Errorf("%02d is not a valid second", sec)
	}

	return time.Date(year, month, day, hour, min, sec, nsec, loc), nil
}

func checkMinSec(val int) bool {
	return val >= 0 && val <= 59
}

func checkHour(val int) bool {
	return val >= 0 && val <= 23
}

func checkMonth(val int) bool {
	return val >= 1 && val <= 12
}

// checkDay returns whether the given day is valid in the given month.  Note that it is not aware of
// leap years, so will allow Feb 29th every year.
func checkDay(month time.Month, day int) bool {
	monthDays := map[time.Month]int{
		time.January:   31,
		time.February:  29,
		time.March:     31,
		time.April:     30,
		time.May:       31,
		time.June:      30,
		time.July:      31,
		time.August:    31,
		time.September: 30,
		time.October:   31,
		time.November:  30,
		time.December:  31,
	}
	max, ok := monthDays[month]
	return ok && day > 0 && day <= max
}

// beginsOffset tells you whether the token is a valid first token for a timezone offset.
func beginsOffset(tok token) bool {
	return tok == DASH || tok == PLUS || tok == Z
}

// parseInt takes a string and casts it to an integer.  It assumes that the string has already been
// validated as having only digit characters inside.  If this assumption is violated, it will panic.
func parseInt(in string) int {
	out, err := strconv.Atoi(in)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// parseDecimal takes in a string like "0234" and returns it as the decimal portion of a float, like 0.0234.
// It assumes that the in string has already been validated as having only digits (so will not error
// on strconv.Atoi), and will panic if that assumption is violated.
func parseDecimal(in string) float64 {
	return float64(parseInt(in)) / math.Pow(10.0, float64(len(in)))
}
