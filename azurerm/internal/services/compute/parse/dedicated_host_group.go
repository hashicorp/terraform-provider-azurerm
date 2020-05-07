package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostGroupId struct {
	ResourceGroup string
	Name          string
}

func DedicatedHostGroupID(input string) (*DedicatedHostGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Dedicated Host Group ID %q: %+v", input, err)
	}

	group := DedicatedHostGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if group.Name, err = id.PopSegment("hostGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
