package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IpGroupResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseIpGroupID(input string) (*IpGroupResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse IP Group ID %q: %+v", input, err)
	}

	ipGroup := IpGroupResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if ipGroup.Name, err = id.PopSegment("ipGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &ipGroup, nil
}
