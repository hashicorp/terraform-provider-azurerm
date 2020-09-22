package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageGalleryId struct {
	ResourceGroup string
	Name          string
}

func NewSharedImageGalleryId(resourceGroup, name string) SharedImageGalleryId {
	return SharedImageGalleryId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id SharedImageGalleryId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func SharedImageGalleryID(input string) (*SharedImageGalleryId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Shared Image Gallery ID %q: %+v", input, err)
	}

	gallery := SharedImageGalleryId{
		ResourceGroup: id.ResourceGroup,
	}

	if gallery.Name, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &gallery, nil
}
