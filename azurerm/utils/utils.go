package utils

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform/helper/schema"
)

func Bool(input bool) *bool {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func String(input string) *string {
	return &input
}

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive) and is divisible by a given number
func IntBetweenDivisibleBy(min, max, divisor int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		if math.Mod(float64(v), float64(divisor)) != 0 {
			es = append(es, fmt.Errorf("expected %s to be divisible by %d", k, divisor))
			return
		}

		return
	}
}
