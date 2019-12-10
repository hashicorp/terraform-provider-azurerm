package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseVirtualHubID(input string) (*VirtualHubResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Hub ID %q: %+v", input, err)
	}

	virtualHub := VirtualHubResourceID{
		Base: *id,
		Name: id.Path["virtualHubs"],
	}

	if virtualHub.Name == "" {
		return nil, fmt.Errorf("ID was missing the `virtualHubs` element")
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
