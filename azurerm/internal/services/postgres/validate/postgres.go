package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
)

func PostgresServerServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z][-0-9a-z]{1,61}[0-9a-z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-'. And must be at least 3 characters and at most 63 characters", k))
	}

	return nil, nil
}

func PostgresDatabaseCollation(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	matched, _ := regexp.MatchString(`^[-A-Za-z0-9_. ]+$`, v)

	if !matched {
		errors = append(errors, fmt.Errorf("%s contains invalid characters, only alphanumeric, underscore, space or hyphen characters are supported, got %s", k, v))
		return
	}

	return warnings, errors
}

func PostgresServerServerID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.PostgresServerServerID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Postgres Server resource id: %v", k, err))
	}

	return warnings, errors
}
