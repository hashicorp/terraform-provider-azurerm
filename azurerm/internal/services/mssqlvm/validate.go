package mssqlvm

import (
	"fmt"
	"regexp"
)

func ValidateVMResourceId(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^/subscriptions/[-\da-zA-Z]/resourcegroups/[-_\da-zA-Z]/providers/Microsoft.Compute/virtualMachines/[-_\da-zA-Z]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must start with '/subscriptions/', invalid Virtual Machine resource ID, ", k))
	}

	return warnings, errors
}

