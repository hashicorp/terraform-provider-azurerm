package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

func VMID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.VirtualMachineID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a virtual machine id: %v", k, err))
	}

	return warnings, errors
}

func MsSqlVMLoginUserName(i interface{}, k string) (warnings []string, errors []error) {
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
