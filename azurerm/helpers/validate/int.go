package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
