package validate

import (
	"fmt"
)

func StringNotEmpty(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q is an empty string", k))
	}

	return
}
