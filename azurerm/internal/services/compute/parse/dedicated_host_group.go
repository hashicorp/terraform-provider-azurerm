package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	HostGroupName  string
}

func NewDedicatedHostGroupID(subscriptionId, resourceGroup, hostGroupName string) DedicatedHostGroupId {
	return DedicatedHostGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		HostGroupName:  hostGroupName,
	}
}

func (id DedicatedHostGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostGroupName)
}

func DedicatedHostGroupID(input string) (*DedicatedHostGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DedicatedHostGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.HostGroupName, err = id.PopSegment("hostGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
