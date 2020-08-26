package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualWanId struct {
	ResourceGroup string
	Name          string
}

func VirtualWanID(input string) (*VirtualWanId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Wan ID %q: %+v", input, err)
	}

	virtualWan := VirtualWanId{
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
