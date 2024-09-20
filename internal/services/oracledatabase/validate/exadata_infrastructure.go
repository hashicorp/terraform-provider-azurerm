package validate

import (
	"fmt"
	"slices"
)

func ComputeCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 2 || v > 32 {
		errors = append(errors, fmt.Errorf("the compute count must be between %d and %d", 2, 32))
		return
	}

	return
}

func StorageCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 3 || v > 64 {
		errors = append(errors, fmt.Errorf("the storage count must be between %d and %d", 3, 64))
		return
	}

	return
}

// MaintenanceWindow validation

func CustomActionTimeoutInMins(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 15 || v > 120 {
		errors = append(errors, fmt.Errorf("the custom action timeout in minutes must be between %d and %d", 15, 120))
		return
	}

	return
}

func LeadTimeInWeeks(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 1 || v > 4 {
		errors = append(errors, fmt.Errorf("the lead time in weeks must be between %d and %d", 1, 4))
		return
	}

	return
}

func Month(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.([]string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be list of strings", k))
		return
	}

	validMonth := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	for _, month := range v {
		if !slices.Contains(validMonth, month) {
			errors = append(errors, fmt.Errorf("month must be %v", validMonth))
			return
		}
	}

	return
}
