package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/parse"
)

// HybridConnectionID validates that the specified ID is a valid Relay Hybrid Connection ID
func HybridConnectionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
	}

	if _, err := parse.HybridConnectionID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", v, err))
	}

	return warnings, errors
}
