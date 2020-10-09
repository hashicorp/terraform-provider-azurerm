package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// FloatInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and matches the value of an element in the valid slice
//
func FloatInSlice(valid []float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(float64)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be float", i))
			return warnings, errors
		}

		for _, validFloat := range valid {
			if v == validFloat {
				return warnings, errors
			}
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %f", k, valid, v))
		return warnings, errors
	}
}
