package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewDedicatedHostGroupId(subscriptionId, resourceGroup, name string) DedicatedHostGroupId {
	return DedicatedHostGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id DedicatedHostGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func DedicatedHostGroupID(input string) (*DedicatedHostGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Dedicated Host Group ID %q: %+v", input, err)
	}

	group := DedicatedHostGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if group.Name, err = id.PopSegment("hostGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
