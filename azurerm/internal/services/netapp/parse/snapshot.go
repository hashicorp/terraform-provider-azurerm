package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SnapshotId struct {
	SubscriptionId    string
	ResourceGroup     string
	NetAppAccountName string
	CapacityPoolName  string
	VolumeName        string
	Name              string
}

func NewSnapshotID(subscriptionId, resourceGroup, netAppAccountName, capacityPoolName, volumeName, name string) SnapshotId {
	return SnapshotId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		NetAppAccountName: netAppAccountName,
		CapacityPoolName:  capacityPoolName,
		VolumeName:        volumeName,
		Name:              name,
	}
}

func (id SnapshotId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Net App Account Name %q", id.NetAppAccountName),
		fmt.Sprintf("Capacity Pool Name %q", id.CapacityPoolName),
		fmt.Sprintf("Volume Name %q", id.VolumeName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id SnapshotId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/capacityPools/%s/volumes/%s/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.Name)
}

// SnapshotID parses a Snapshot ID into an SnapshotId struct
func SnapshotID(input string) (*SnapshotId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SnapshotId{
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
	if resourceId.VolumeName, err = id.PopSegment("volumes"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("snapshots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
