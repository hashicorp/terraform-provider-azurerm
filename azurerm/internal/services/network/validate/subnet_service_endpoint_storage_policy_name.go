package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func SubnetServiceEndpointStoragePolicyName(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_]([\w.\-]{0,78}[\w])?$`), "The name can be up to 80 characters long. It must begin with a alphnum character, and it must end with a alphnum character or with '_'. The name may contain alphnum characters or '.', '-', '_'.")(i, k)
}

func SubnetServiceEndpointStoragePolicyDefinitionName(i interface{}, k string) (warnings []string, errors []error) {
	// Same rule as policy
	return SubnetServiceEndpointStoragePolicyName(i, k)
}
