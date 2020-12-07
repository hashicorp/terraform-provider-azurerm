package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageId struct {
	SubscriptionId string
	ResourceGroup  string
	GalleryName    string
	ImageName      string
}

func NewSharedImageID(subscriptionId, resourceGroup, galleryName, imageName string) SharedImageId {
	return SharedImageId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		GalleryName:    galleryName,
		ImageName:      imageName,
	}
}

func (id SharedImageId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Image Name %q", id.ImageName),
	}
	return strings.Join(segments, " / ")
}

func (id SharedImageId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ImageName)
}

// SharedImageID parses a SharedImage ID into an SharedImageId struct
func SharedImageID(input string) (*SharedImageId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SharedImageId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.GalleryName, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}
	if resourceId.ImageName, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
