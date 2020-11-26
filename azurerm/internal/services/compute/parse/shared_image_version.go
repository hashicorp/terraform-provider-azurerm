package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageVersionId struct {
	SubscriptionId string
	ResourceGroup  string
	GalleryName    string
	ImageName      string
	VersionName    string
}

func NewSharedImageVersionID(subscriptionId, resourceGroup, galleryName, imageName, versionName string) SharedImageVersionId {
	return SharedImageVersionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		GalleryName:    galleryName,
		ImageName:      imageName,
		VersionName:    versionName,
	}
}

func (id SharedImageVersionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName)
}

// SharedImageVersionID parses a SharedImageVersion ID into an SharedImageVersionId struct
func SharedImageVersionID(input string) (*SharedImageVersionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SharedImageVersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.GalleryName, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}
	if resourceId.ImageName, err = id.PopSegment("images"); err != nil {
		return nil, err
	}
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
