package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func StatusCodeRange(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to contain a single '-', got %v", k, i))
		return warnings, errors
	}

	if _, err := strconv.Atoi(parts[0]); err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the left of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	if _, err := strconv.Atoi(parts[1]); err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the right of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	return warnings, errors
}
