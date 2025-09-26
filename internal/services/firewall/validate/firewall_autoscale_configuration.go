package validate

import "fmt"

func AutoscaleConfiguration(d map[string]interface{}) (warnings []string, errors []error) {
	minPresent := false
	maxPresent := false
	min := 0
	max := 0

	if _, ok := d["minCapacity"]; ok {
		minPresent = true
		min = d["minCapacity"].(int)
	}
	if _, ok := d["maxCapacity"]; ok {
		maxPresent = true
		max = d["maxCapacity"].(int)
	}
	if minPresent && maxPresent {
		if min > max {
			errors = append(errors, fmt.Errorf("minCapacity cannot be greater than maxCapacity"))
		}
		if max-min < 2 || max-min != 0 {
			errors = append(errors, fmt.Errorf("maxCapacity must be at least 2 greater than minCapacity or be exactly 0"))
		}
	}
	return warnings, errors
}

func MinCapacity(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 2 || value > 50 {
		errors = append(errors, fmt.Errorf("minCapacity must be between 2 and 50"))
	}
	return warnings, errors
}

func MaxCapacity(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 2 || value > 50 {
		errors = append(errors, fmt.Errorf("maxCapacity must be between 2 and 50"))
	}
	return warnings, errors
}
