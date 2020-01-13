package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualWanResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseVirtualWanID(input string) (*VirtualWanResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Wan ID %q: %+v", input, err)
	}

	virtualWan := VirtualWanResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualWan.Name, err = id.PopSegment("virtualWans"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualWan, nil
}

func ValidateVirtualWanID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := ParseVirtualWanID(v); err != nil {
		return nil, []error{err}
	}

	return nil, nil
}
