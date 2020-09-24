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

func FrontDoorRulesEngineName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^([A-Za-z])([A-Za-z\d]{0,255})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must begin with a letter, be between 1 and 256 characters in length and may contain only letters or numbers.`, k))
	}
	
	return nil, nil
}

func FrontDoorWAFName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{0,127})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 128 characters in length, must begin with a letter and may only contain letters and numbers.`, k))
	}

	return nil, nil
}

func FrontDoorBackendPoolRoutingRuleName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{1,88})([\da-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, nil
}

func FrontdoorCustomBlockResponseBody(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q contains invalid characters, %q must contain only alphanumeric and equals sign characters.`, k, k))
	}

	return nil, nil
}
