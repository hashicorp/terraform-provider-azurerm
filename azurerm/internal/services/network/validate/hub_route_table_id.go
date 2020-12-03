package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func HubRouteTableID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := parse.HubRouteTableID(v); err != nil {
		return nil, []error{err}
	}

	return nil, nil
}
