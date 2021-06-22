package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

// Your server name can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.
func ValidateMsSqlServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.", k))
	}

	return nil, nil
}

// Your database name can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters.
func ValidateMsSqlDatabaseName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[^<>*%&:\\\/?]{0,127}[^\s.<>*%&:\\\/?]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q can't end with '.' or ' ', can't contain '<,>,*,%%,&,:,\,/,?' or control characters, and can't have more than 128 characters.`, k))
	}

	return nil, nil
}

func ValidateMsSqlFailoverGroupName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-'.", k))
	}

	return nil, nil
}

// Following characters and any control characters are not allowed for resource name '%,&,\\\\,?,/'.\"
// The name can not end with characters: '. '
// TODO: unsure about length, was able to deploy one at 120
func ValidateMsSqlElasticPoolName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[^&%\\\/?]{0,127}[^\s.&%\\\/?]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q can't end with '.' or ' ', can't contain '%%,&,\,/,?' or control characters, and can't have more than 128 characters.`, k))
	}

	return nil, nil
}

// Job Agent name must not contain any of ?<>*%&:\/? and must not end with a space or .
func ValidateMsSqlJobAgentName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[^?<>*%&:\/?]{0,127}[^?<>*%&:\/?. ]$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q must not contain any of ?<>*%%&:\\/?, must not end with a space or a period and can't have more than 128 characters", k))
	}

	return nil, nil
}
