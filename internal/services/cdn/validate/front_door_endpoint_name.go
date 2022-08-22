package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorEndpointName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[\da-zA-Z][-\da-zA-Z]{0,44}[\da-zA-Z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 2 and 46 characters in length, begin with a letter or number, end with a letter or number and may contain only letters, numbers and hyphens, got %q`, k, i))
	}

	return nil, nil
}
