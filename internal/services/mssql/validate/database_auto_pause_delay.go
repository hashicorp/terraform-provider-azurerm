package validate

import "fmt"

func DatabaseAutoPauseDelay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
		return warnings, errors
	}
	min := 60
	max := 10080
	if (v < min || v > max) && v%10 != 0 && v != -1 {
		errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d) and divisible by 10 or -1, got %d", k, min, max, v))
		return warnings, errors
	}

	return warnings, errors
}
