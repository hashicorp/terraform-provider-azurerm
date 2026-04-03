// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "fmt"

func DiskSizeGB(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 32767 {
		errors = append(errors, fmt.Errorf(
			"%s can only be between 0 and 32767", k))
	}
	return warnings, errors
}
