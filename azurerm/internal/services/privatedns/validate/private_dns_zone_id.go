package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func PrivateDnsZoneID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// TODO: switch this to use an ID Parser
	if id, err := azure.ParseAzureResourceID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
	} else if _, err = id.PopSegment("privateDnsZones"); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a private dns zone resource id: %v", k, err))
	}

	return warnings, errors
}
