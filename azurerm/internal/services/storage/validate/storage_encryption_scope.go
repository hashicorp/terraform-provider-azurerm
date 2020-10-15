package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func StorageEncryptionScopeName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[0-9a-zA-Z]{4,63}$").MatchString(input) {
		errors = append(errors, fmt.Errorf("storage encryption scope name %q must be alphanumeric, and between 4 to 63 characters", input))
	}

	return warnings, errors
}

func KeyVaultChildId(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := azure.ParseKeyVaultChildID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Key Vault Child resource id: %v", k, err))
	}

	return warnings, errors
}
