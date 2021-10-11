package validate

import "fmt"

func DiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 128 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` must be 128 or greater"))
	}
	return warnings, errors
}
