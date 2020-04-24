package validate

import (
	"fmt"
	"time"
)

func Duration(i interface{}, k string) (warnings []string, errors []error) {
	value := i.(string)
	duration, err := time.ParseDuration(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as a duration: %s", k, err))
	}
	if duration < 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than zero", k))
	}
	return warnings, errors
}
