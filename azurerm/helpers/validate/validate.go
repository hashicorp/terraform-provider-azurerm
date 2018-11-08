package validate

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func IntBetweenAndNot(min, max, not int) schema.SchemaValidateFunc {
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

		if v == not {
			errors = append(errors, fmt.Errorf("expected %s to not be %d, got %d", k, not, v))
			return
		}

		return
	}
}
