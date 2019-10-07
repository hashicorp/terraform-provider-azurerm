package netapp

import (
	"fmt"
	"net"

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

func ValidateActiveDirectoryDNSName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	result := net.ParseIP(v)
	if result.To4() == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 IP address", v))
	}

	return nil, errors
}
