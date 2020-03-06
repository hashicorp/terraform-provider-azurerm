package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
)

func DomainRegistrationID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.DomainRegistrationID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func DomainRegistrationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// TODO - list of supported TLDs to check against

	parts := strings.Split(value, ".")

	if matched := regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{0,61}[0-9a-zA-Z]$`).Match([]byte(parts[0])); !matched {
		errors = append(errors, fmt.Errorf("%q domain part may only contain alphanumeric characters and dashes up to 63 characters in length, and must start and end in an alphanumeric", k))
	}

	return warnings, errors
}
