package azure

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

//Frontdoor name must begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.
func ValidateFrontDoorName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[\da-zA-Z]([-\da-z]{4,61})[\da-zA-Z]?$`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, errors
}
