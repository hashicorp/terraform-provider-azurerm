package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
)

func KubernetesClusterID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.KubernetesClusterID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Resource Id: %v", v, err))
	}

	return warnings, errors
}
