package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func IsGatewaySubnet(i interface{}, k string) (warnings []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := parse.SubnetIDInsensitively(value)
	if err != nil {
		errors = append(errors, err)
		return
	}

	if !strings.EqualFold(id.Name, "GatewaySubnet") {
		errors = append(errors, fmt.Errorf("expected %s to reference a gateway subnet with name GatewaySubnet", k))
	}

	return warnings, errors
}
