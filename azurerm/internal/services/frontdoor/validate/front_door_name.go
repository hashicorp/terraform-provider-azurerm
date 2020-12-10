package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func FrontDoorName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{3,61})([\da-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 5 and 63 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}
	return nil, nil
}
