package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseVirtualHubID(input string) (*VirtualHubResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Hub ID %q: %+v", input, err)
	}

	virtualHub := VirtualHubResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHub.Name, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHub, nil
}

func ValidateVirtualHubID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := ParseVirtualHubID(v); err != nil {
		return nil, []error{err}
	}

	return nil, nil
}
