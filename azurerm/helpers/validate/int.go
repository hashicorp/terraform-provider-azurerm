package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// deprecated: use validation.All(validation.IntBetween, validation.IntNotInSlice)
func IntBetweenAndNot(min, max, not int) schema.SchemaValidateFunc {
	return validation.All(validation.IntBetween(min, max), validation.IntNotInSlice([]int{not}))
}

func IntBetweenAndNotInRange(min, max, rangeMin, rangeMax int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be int", k))
			return
		}
		if v < min || v > max {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		if v >= rangeMin && v <= rangeMax {
			errors = append(errors, fmt.Errorf("expected %s to not be in the range (%d - %d), got %d", k, rangeMin, rangeMax, v))
			return
		}

		return
	}
}

// deprecated: use validation.All(validation.IntBetween, validation.IntDivisibleBy)
func IntBetweenAndDivisibleBy(min, max, divisor int) schema.SchemaValidateFunc { // nolint: unparam
	return validation.All(validation.IntBetween(min, max), validation.IntDivisibleBy(divisor))
}

// deprecated: use validation.IntDivisibleBy
func IntDivisibleBy(divisor int) schema.SchemaValidateFunc { // nolint: unparam
	return validation.IntDivisibleBy(divisor)
}

// deprecated: use validation.IntInSlice
func IntInSlice(valid []int) schema.SchemaValidateFunc {
	return validation.IntInSlice(valid)
}
