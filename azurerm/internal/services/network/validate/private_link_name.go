package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func PrivateLinkName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules per the Nat Gateway service team are (Friday, October 18, 2019 4:20 PM):
	// 1. Must not be empty.
	// 2. Must be between 1 and 80 characters.
	// 3. The attribute must:
	//    a) begin with a letter or number
	//    b) end with a letter, number or underscore
	//    c) may contain only letters, numbers, underscores, periods, or hyphens.

	if len(v) == 1 {
		if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z\d])`); !m {
			errors = append(errors, fmt.Errorf("%s must begin with a letter or number", k))
		}
	} else if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z\d])([a-zA-Z\d-\_\.]{0,78})([a-zA-Z\d\_])$`); !m {
		errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, periods, hyphens or underscores", k))
	}

	return nil, errors
}
