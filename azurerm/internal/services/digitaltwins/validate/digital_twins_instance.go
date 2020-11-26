package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
)

func DigitaltwinsInstanceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 3 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 3, v))
		return
	}

	if len(v) > 63 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 63, v))
		return
	}

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter or number, end with a letter or number and contain only letters, numbers, and hyphens, got %v", k, v))
		return
	}
	return
}

func DigitaltwinsInstanceID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.DigitalTwinsInstanceID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Digital Twins Instance resource id: %v", k, err))
	}

	return warnings, errors
}
