package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageSyncGroupId struct {
	SubscriptionId         string
	ResourceGroup          string
	StorageSyncServiceName string
	SyncGroupName          string
}

func NewStorageSyncGroupID(subscriptionId, resourceGroup, storageSyncServiceName, syncGroupName string) StorageSyncGroupId {
	return StorageSyncGroupId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		StorageSyncServiceName: storageSyncServiceName,
		SyncGroupName:          syncGroupName,
	}
}

func (id StorageSyncGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Storage Sync Service Name %q", id.StorageSyncServiceName),
		fmt.Sprintf("Sync Group Name %q", id.SyncGroupName),
	}
	return strings.Join(segments, " / ")
}

func (id StorageSyncGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s/syncGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName)
}

// StorageSyncGroupID parses a StorageSyncGroup ID into an StorageSyncGroupId struct
func StorageSyncGroupID(input string) (*StorageSyncGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageSyncGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageSyncServiceName, err = id.PopSegment("storageSyncServices"); err != nil {
		return nil, err
	}
	if resourceId.SyncGroupName, err = id.PopSegment("syncGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
