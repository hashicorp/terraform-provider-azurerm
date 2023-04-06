package validate

import (
	"fmt"
)

func PartitionCount(v interface{}, k string) (warnings []string, errors []error) {
	partitionCount := v.(int)

	switch partitionCount {
	case 5, 7, 8, 9, 10, 11:
		errors = append(errors, fmt.Errorf("%q must be 1, 2, 3, 4, 6, or 12, got %d", k, partitionCount))
	}

	if partitionCount > 12 {
		errors = append(errors, fmt.Errorf("%q must be 1, 2, 3, 4, 6, or 12, got %d", k, partitionCount))
	}

	return warnings, errors
}
