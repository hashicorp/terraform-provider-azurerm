package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageVersionId struct {
	SubscriptionId string
	ResourceGroup  string
	Gallery        string
	ImageName      string
	Version        string
}

func NewSharedImageVersionId(id SharedImageId, name string) SharedImageVersionId {
	return SharedImageVersionId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		Gallery:        id.Gallery,
		ImageName:      id.Name,
		Version:        name,
	}
}

func (id SharedImageVersionId) ID(_ string) string {
	galleryId := NewSharedImageGalleryId(id.SubscriptionId, id.ResourceGroup, id.Gallery)
	base := NewSharedImageId(galleryId, id.ImageName).ID("")
	return fmt.Sprintf("%s/versions/%s", base, id.Version)
}

func SharedImageVersionID(input string) (*SharedImageVersionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Shared Image Version ID %q: %+v", input, err)
	}

	set := SharedImageVersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if set.Gallery, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if set.ImageName, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if set.Version, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
