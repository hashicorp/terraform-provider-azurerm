package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkID struct {
	ResourceGroup string
	Name          string
}

func ParseVirtualNetworkID(input string) (*VirtualNetworkID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Subnet ID %q: %+v", input, err)
	}

	vnet := VirtualNetworkID{
		ResourceGroup: id.ResourceGroup,
	}

	if vnet.Name, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vnet, nil
}

// ValidateVirtualNetworkID validates that the specified ID is a valid App Service ID
func ValidateVirtualNetworkID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseVirtualNetworkID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
