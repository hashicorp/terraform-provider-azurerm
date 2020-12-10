package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func FrontDoorWAFName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{0,127})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 128 characters in length, must begin with a letter and may only contain letters and numbers.`, k))
	}

	return nil, nil
}
