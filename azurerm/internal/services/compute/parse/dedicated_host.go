package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id DedicatedHostId) String() string {
	segments := []string{
		fmt.Sprintf("Host Name %q", id.HostName),
		fmt.Sprintf("Host Group Name %q", id.HostGroupName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dedicated Host", segmentsStr)
}

func (id DedicatedHostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s/hosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostGroupName, id.HostName)
}

// DedicatedHostID parses a DedicatedHost ID into an DedicatedHostId struct
func DedicatedHostID(input string) (*DedicatedHostId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DedicatedHostId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
