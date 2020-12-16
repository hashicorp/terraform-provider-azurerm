package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DiskAccessId struct {
	SubscriptionId string
	ResourceGroup  string
	DiskAccessName string
}

func NewDiskAccessID(subscriptionId, resourceGroup, diskAccessName string) DiskAccessId {
	return DiskAccessId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DiskAccessName: diskAccessName,
	}
}

func (id DiskAccessId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Disk Access Name %q", id.DiskAccessName),
	}
	return strings.Join(segments, " / ")
}

func (id DiskAccessId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskAccesses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DiskAccessName)
}

// DiskAccessID parses a DiskAccess ID into an DiskAccessID struct
func DiskAccessID(input string) (*DiskAccessId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DiskAccessId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DiskAccessName, err = id.PopSegment("diskAccesses"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
