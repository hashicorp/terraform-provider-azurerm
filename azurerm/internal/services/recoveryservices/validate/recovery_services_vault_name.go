package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func RecoveryServicesVaultName(v interface{}, k string) (warnings []string, errors []error) {
	regexpValidator := validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,49}$"),
		"Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens.",
	)
	return regexpValidator(v, k)
}
