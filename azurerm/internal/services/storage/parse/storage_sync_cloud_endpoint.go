package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageSyncCloudEndpointId struct {
	SubscriptionId         string
	ResourceGroup          string
	StorageSyncServiceName string
	SyncGroupName          string
	CloudEndpointName      string
}

func NewStorageSyncCloudEndpointID(subscriptionId, resourceGroup, storageSyncServiceName, syncGroupName, cloudEndpointName string) StorageSyncCloudEndpointId {
	return StorageSyncCloudEndpointId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		StorageSyncServiceName: storageSyncServiceName,
		SyncGroupName:          syncGroupName,
		CloudEndpointName:      cloudEndpointName,
	}
}

func (id StorageSyncCloudEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Cloud Endpoint Name %q", id.CloudEndpointName),
		fmt.Sprintf("Sync Group Name %q", id.SyncGroupName),
		fmt.Sprintf("Storage Sync Service Name %q", id.StorageSyncServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Sync Cloud Endpoint", segmentsStr)
}

func (id StorageSyncCloudEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s/syncGroups/%s/cloudEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName, id.CloudEndpointName)
}

// StorageSyncCloudEndpointID parses a StorageSyncCloudEndpoint ID into an StorageSyncCloudEndpointId struct
func StorageSyncCloudEndpointID(input string) (*StorageSyncCloudEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageSyncCloudEndpointId{
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
	if resourceId.CloudEndpointName, err = id.PopSegment("cloudEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
