package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func VirtualHubRouteTableID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := parse.VirtualHubRouteTableID(v); err != nil {
		return nil, []error{err}
	}

	return nil, nil
}

func VirtualHubRouteTableName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[^<>%&:?/+]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must not contain characters from %q", k, "<>&:?/+%"))
	}

	return warnings, errors
}
