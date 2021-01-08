package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VolumeId struct {
	SubscriptionId    string
	ResourceGroup     string
	NetAppAccountName string
	CapacityPoolName  string
	Name              string
}

func NewVolumeID(subscriptionId, resourceGroup, netAppAccountName, capacityPoolName, name string) VolumeId {
	return VolumeId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		NetAppAccountName: netAppAccountName,
		CapacityPoolName:  capacityPoolName,
		Name:              name,
	}
}

func (id VolumeId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Capacity Pool Name %q", id.CapacityPoolName),
		fmt.Sprintf("Net App Account Name %q", id.NetAppAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Volume", segmentsStr)
}

func (id VolumeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/capacityPools/%s/volumes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name)
}

// VolumeID parses a Volume ID into an VolumeId struct
func VolumeID(input string) (*VolumeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VolumeId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}
	if resourceId.CapacityPoolName, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("volumes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
