package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func FirewallPolicyPolicyName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_][\w-.]*[\w]$`),
		"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.")
}

func FirewallPolicyPolicyID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.FirewallPolicyPolicyID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
