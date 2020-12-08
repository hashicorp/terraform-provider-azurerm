package validate

import "fmt"

func ClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if len(v) == 0 || len(v) > 260 {
		errors = append(errors, fmt.Errorf("%s cannot be empty and must not exceed 260 characters", k))
		return warnings, errors
	}

	return warnings, errors
}
