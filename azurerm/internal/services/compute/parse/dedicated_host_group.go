package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostGroupId struct {
	ResourceGroup string
	Name          string
}

func NewDedicatedHostGroupId(resourceGroup, name string) DedicatedHostGroupId {
	return DedicatedHostGroupId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id DedicatedHostGroupId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DedicatedHostGroupID(input string) (*DedicatedHostGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Dedicated Host Group ID %q: %+v", input, err)
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
