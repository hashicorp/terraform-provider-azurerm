package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageContainerResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	BlobServiceName    string
	ContainerName      string
}

func NewStorageContainerResourceManagerID(subscriptionId, resourceGroup, storageAccountName, blobServiceName, containerName string) StorageContainerResourceManagerId {
	return StorageContainerResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		BlobServiceName:    blobServiceName,
		ContainerName:      containerName,
	}
}

func (id StorageContainerResourceManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Container Name %q", id.ContainerName),
		fmt.Sprintf("Blob Service Name %q", id.BlobServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Container Resource Manager", segmentsStr)
}

func (id StorageContainerResourceManagerId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/%s/containers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.BlobServiceName, id.ContainerName)
}

// StorageContainerResourceManagerID parses a StorageContainerResourceManager ID into an StorageContainerResourceManagerId struct
func StorageContainerResourceManagerID(input string) (*StorageContainerResourceManagerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageContainerResourceManagerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.BlobServiceName, err = id.PopSegment("blobServices"); err != nil {
		return nil, err
	}
	if resourceId.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
