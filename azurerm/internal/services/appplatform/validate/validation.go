package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func ValidateSpringCloudServiceName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules are :
	// 1. can contain only lowercase letters, numbers and hyphens.
	// 2. The first character must be a letter.
	// 3. The last character must be a letter or number
	// 3. The value must be between 4 and 32 characters long

	if len(v) < 4 || len(v) > 32 {
		errors = append(errors, fmt.Errorf("%s must be between 4 and 32 characters long", k))
	} else {
		if m, _ := validate.RegExHelper(i, k, `^([a-z])([a-z\d-]{0,30})([a-z\d])$`); !m {
			errors = append(errors, fmt.Errorf("%s must begin with a letter, end with a letter or number, contain only lowercase letters, numbers and hyphens.", k))
		}
	}

	return nil, errors
}

func ValidateConfigServerURI(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// the config server URI should be started with http://, https://, git@, or ssh://
	if !strings.HasPrefix(v, "http://") &&
		!strings.HasPrefix(v, "https://") &&
		!strings.HasPrefix(v, "git@") &&
		!strings.HasPrefix(v, "ssh://") {
		errors = append(errors, fmt.Errorf("%s should be started with http://, https://, git@, or ssh://", k))
	}
	return nil, errors
}
