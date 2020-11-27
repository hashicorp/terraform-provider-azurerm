package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IpGroupId struct {
	ResourceGroup string
	Name          string
}

func IpGroupID(input string) (*IpGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse IP Group ID %q: %+v", input, err)
	}

	ipGroup := IpGroupId{
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
