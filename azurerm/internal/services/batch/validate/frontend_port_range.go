package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func FrontendPortRange(i interface{}, k string) (warnings []string, errors []error) {
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

	startPort, err := strconv.Atoi(parts[0])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the left of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	endPort, err := strconv.Atoi(parts[1])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the right of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	if !validPortNumber(startPort) || !validPortNumber(endPort) {
		errors = append(errors, fmt.Errorf("expect values range between 1 and 65534 except ports from `50000` to `55000`, got %v: %v", k, i))
		return warnings, errors
	}
	if endPort-startPort < 100 {
		errors = append(errors, fmt.Errorf("values must be a range of at least 100, got %v: %v", k, i))
		return warnings, errors
	}

	return warnings, errors
}

func validPortNumber(port int) bool {
	return 1 <= port && port < 50000 || 55000 < port && port <= 65535
}
