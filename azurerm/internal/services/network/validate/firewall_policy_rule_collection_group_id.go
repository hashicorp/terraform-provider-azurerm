package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func FirewallPolicyRuleCollectionGroupID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.FirewallPolicyRuleCollectionGroupID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
