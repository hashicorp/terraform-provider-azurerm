package validate

import (
	"fmt"
	"strconv"
)

// MaxBidPrice validates `max_bid_price` against the following rules:
// * must be either -1 or positive number
// * must have up to 5 digits after the radix point
func MaxBidPrice(value interface{}, key string) (warns []string, errors []error) {
	v := value.(float64)

	if v < -1 || v > -1 && v <= 0 {
		errors = append(errors, fmt.Errorf("%q must be between 0 (exclusive) and `math.MaxFloat64` (inclusive) or -1 (special value), got %f", key, v))
	}

	strV := fmt.Sprintf("%.5f", v)
	parsedFloat, _ := strconv.ParseFloat(strV, 64)

	if parsedFloat != v {
		errors = append(errors, fmt.Errorf("%q can only include up to 5 digits after the radix point, got %g", key, v))
	}

	return
}
