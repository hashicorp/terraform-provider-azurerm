package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostId struct {
	SubscriptionId string
	ResourceGroup  string
	HostGroupName  string
	HostName       string
}

func NewDedicatedHostID(subscriptionId, resourceGroup, hostGroupName, hostName string) DedicatedHostId {
	return DedicatedHostId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		HostGroupName:  hostGroupName,
		HostName:       hostName,
	}
}

func (id DedicatedHostId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s/hosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostGroupName, id.HostName)
}

func DedicatedHostID(input string) (*DedicatedHostId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DedicatedHostId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.HostGroupName, err = id.PopSegment("hostGroups"); err != nil {
		return nil, err
	}
	if resourceId.HostName, err = id.PopSegment("hosts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
