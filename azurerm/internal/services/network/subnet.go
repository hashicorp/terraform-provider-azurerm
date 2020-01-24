package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubnetID struct {
	ResourceGroup      string
	VirtualNetworkName string
	Name               string
}

func ParseSubnetID(input string) (*SubnetID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Subnet ID %q: %+v", input, err)
	}

	subnet := SubnetID{
		ResourceGroup: id.ResourceGroup,
	}

	if subnet.VirtualNetworkName, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}

	if subnet.Name, err = id.PopSegment("subnets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &subnet, nil
}

// ValidateSubnetID validates that the specified ID is a valid App Service ID
func ValidateSubnetID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseSubnetID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
