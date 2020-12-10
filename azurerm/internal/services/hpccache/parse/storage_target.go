package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageTargetId struct {
	SubscriptionId string
	ResourceGroup  string
	CacheName      string
	Name           string
}

func NewStorageTargetID(subscriptionId, resourceGroup, cacheName, name string) StorageTargetId {
	return StorageTargetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		CacheName:      cacheName,
		Name:           name,
	}
}

func (id StorageTargetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Cache Name %q", id.CacheName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Target", segmentsStr)
}

func (id StorageTargetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/caches/%s/storageTargets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CacheName, id.Name)
}

// StorageTargetID parses a StorageTarget ID into an StorageTargetId struct
func StorageTargetID(input string) (*StorageTargetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageTargetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CacheName, err = id.PopSegment("caches"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("storageTargets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
