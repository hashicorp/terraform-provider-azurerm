package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageGalleryId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSharedImageGalleryId(subscriptionId, resourceGroup, name string) SharedImageGalleryId {
	return SharedImageGalleryId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id SharedImageGalleryId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func SharedImageGalleryID(input string) (*SharedImageGalleryId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Shared Image Gallery ID %q: %+v", input, err)
	}

	gallery := SharedImageGalleryId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if gallery.Name, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &gallery, nil
}
