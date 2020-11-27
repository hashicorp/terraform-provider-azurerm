package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

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

func (id VolumeId) ID(_ string) string {
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
