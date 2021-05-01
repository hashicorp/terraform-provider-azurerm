package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func PrivateLinkSubResourceName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if len(strings.TrimSpace(v)) >= 3 {
		if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z0-9])([\w\.-]{1,61})([a-zA-Z0-9])$`); !m {
			errors = append(errors, fmt.Errorf("%s must begin and end with a alphanumeric character, be between 3 and 63 characters in length, only contain letters, numbers, underscores, periods, and dashes", k))
		}
	} else {
		errors = append(errors, fmt.Errorf("%s must be at least 3 character in length", k))
	}

	return nil, errors
}
