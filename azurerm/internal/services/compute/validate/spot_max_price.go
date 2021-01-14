package validate

import (
	"fmt"
)

// SpotMaxPrice validates the price provided is a valid Spot Price for the Compute
// API (and downstream API's which use this like AKS)
func SpotMaxPrice(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float", k))
		return
	}

	// either -1 (the current VM price)
	if v == -1.0 {
		return
	}

	// at least 0.00001
	if v < 0.00001 {
		errors = append(errors, fmt.Errorf("expected %q to be > 0.00001 but got %.5f", k, v))
		return
	}

	return
}
