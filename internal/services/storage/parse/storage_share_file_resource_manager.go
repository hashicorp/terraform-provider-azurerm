package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageShareFileResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	FileServiceName    string
	ShareName          string
}

func NewStorageShareFileResourceManagerID(subscriptionId, resourceGroup, storageAccountName, fileServiceName, shareName string) StorageShareFileResourceManagerId {
	return StorageShareFileResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		FileServiceName:    fileServiceName,
		ShareName:          shareName,
	}
}

func (id StorageShareFileResourceManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Share Name %q", id.ShareName),
		fmt.Sprintf("File Service Name %q", id.FileServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Share File Resource Manager", segmentsStr)
}

func (id StorageShareFileResourceManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/shares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileServiceName, id.ShareName)
}

// StorageShareFileResourceManagerID parses a StorageShareFileResourceManager ID into an StorageShareFileResourceManagerId struct
func StorageShareFileResourceManagerID(input string) (*StorageShareFileResourceManagerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageShareFileResourceManagerId{
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
	if resourceId.FileServiceName, err = id.PopSegment("fileServices"); err != nil {
		return nil, err
	}
	if resourceId.ShareName, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
