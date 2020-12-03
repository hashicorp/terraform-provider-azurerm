package validate

import (
	"fmt"
	"regexp"
)

func SqlVirtualMachineLoginUserName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile(`^[^\\/"\[\]:|<>+=;,?* .]{2,128}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%v cannot contain special characters '\\/\"[]:|<>+=;,?* .'", k))
	}

	return warnings, errors
}
