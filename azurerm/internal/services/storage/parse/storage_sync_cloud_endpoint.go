package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type SyncCloudEndpointId struct {
	Name             string
	StorageSyncName  string
	StorageSyncGroup string
	ResourceGroup    string
}

func SyncCloudEndpointID(input string) (*SyncCloudEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cloudEndpoint := SyncCloudEndpointId{
		ResourceGroup: id.ResourceGroup,
	}

	if cloudEndpoint.StorageSyncName, err = id.PopSegment("storageSyncServices"); err != nil {
		return nil, err
	}

	if cloudEndpoint.StorageSyncGroup, err = id.PopSegment("syncGroups"); err != nil {
		return nil, err
	}

	if cloudEndpoint.Name, err = id.PopSegment("cloudEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cloudEndpoint, nil
}
