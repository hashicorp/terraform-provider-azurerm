package validate

import "fmt"

// ScheduledQueryRulesAlertThreshold checks that a threshold value is between 0 and 10000
// and is a whole number. The azure-sdk-for-go expects this value to be a float64
// but the user validation rules want an integer.
func ScheduledQueryRulesAlertThreshold(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float64", k))
	}

	if v != float64(int64(v)) {
		errors = append(errors, fmt.Errorf("%q must be a whole number", k))
	}

	if v < 0 || v > 10000 {
		errors = append(errors, fmt.Errorf("%q must be between 0 and 10000 (inclusive)", k))
	}

	return warnings, errors
}
