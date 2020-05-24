package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
)

func ManagementGroupName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// portal says: The name can only be an ASCII letter, digit, -, _, (, ), . and have a maximum length constraint of 90
	if matched := regexp.MustCompile(`^[a-zA-Z0-9_().-]{1,90}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s can only consist of ASCII letters, digits, -, _, (, ), . , and cannot exceed the maximum length of 90", k))
	}
	return
}

func ManagementGroupID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ManagementGroupID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a management group id: %v", k, err))
		return
	}

	return
}
