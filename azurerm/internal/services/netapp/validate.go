package netapp

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func ValidateNetAppAccountName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{1,62})([\da-zA-Z]$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must be between 3 and 64 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, errors
}

func ValidateActiveDirectoryDomainName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^(.*)[\da-zA-Z]\.[\da-zA-Z].{1,}$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must end with a letter or number before dot and start with a letter or number after dot.`, k))
	}

	if len(k) > 255 {
		errors = append(errors, fmt.Errorf(`Active Directory Domain Name can not be longer than 255 characters in length, got %d characters`, len(k)))
	}

	return nil, errors
}
