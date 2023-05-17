package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SubnetServiceEndpointStoragePolicyAlias(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^\/services\/`), "The name must start with /services/")(i, k)
}

// func SubnetServiceEndpointStoragePolicyDefinitionName(i interface{}, k string) (warnings []string, errors []error) {
// 	// Same rule as policy
// 	return SubnetServiceEndpointStoragePolicyName(i, k)
// }
