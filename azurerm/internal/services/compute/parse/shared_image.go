package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageId struct {
	ResourceGroup string
	Gallery       string
	Name          string
}

func NewSharedImageId(id SharedImageGalleryId, name string) SharedImageId {
	return SharedImageId{
		ResourceGroup: id.ResourceGroup,
		Gallery:       id.Name,
		Name:          name,
	}
}

func (id SharedImageId) ID(subscriptionId string) string {
	base := NewSharedImageGalleryId(id.ResourceGroup, id.Gallery).ID(subscriptionId)
	return fmt.Sprintf("%s/images/%s", base, id.Name)
}

func SharedImageID(input string) (*SharedImageId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Shared Image ID %q: %+v", input, err)
	}

	image := SharedImageId{
		ResourceGroup: id.ResourceGroup,
	}

	if image.Gallery, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if image.Name, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &image, nil
}
