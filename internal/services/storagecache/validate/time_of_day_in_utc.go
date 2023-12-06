package validate

import (
	"fmt"
	"regexp"
)

func TimeOfDayInUTC(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9]):([0-5][0-9])$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must match the format HHmm where HH is 00-23 and mm is 00-59", k))
	}

	return warnings, errors
}
