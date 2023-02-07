package validate

import (
	"fmt"
)

// FloatInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and matches the value of an element in the valid slice
func FloatInSlice(valid []float64) func(interface{}, string) ([]string, []error) {
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

func FloatInRange(start, end float64) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(float64)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be float", i))
			return warnings, errors
		}

		if v < start || v > end {
			errors = append(errors, fmt.Errorf("expected %s to be in range [%f, %f], got %f", k, start, end, v))
		}

		return warnings, errors
	}
}

func IntegerPositive(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", i))
		return
	}
	if v <= 0 {
		errors = append(errors, fmt.Errorf("expected %s to be positive, got %d", k, v))
		return
	}
	return
}
