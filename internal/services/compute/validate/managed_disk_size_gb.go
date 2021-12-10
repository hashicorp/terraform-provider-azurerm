package validate

import "fmt"

func ManagedDiskSizeGB(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 65536 {
		errors = append(errors, fmt.Errorf("%q can only be between 0 and 65536", k))
	}
	return warnings, errors
}
