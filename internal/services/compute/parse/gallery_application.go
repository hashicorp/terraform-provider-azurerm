package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GalleryApplicationId struct {
	SubscriptionId  string
	ResourceGroup   string
	GalleryName     string
	ApplicationName string
}

func NewGalleryApplicationID(subscriptionId, resourceGroup, galleryName, applicationName string) GalleryApplicationId {
	return GalleryApplicationId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		GalleryName:     galleryName,
		ApplicationName: applicationName,
	}
}

func (id GalleryApplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Application Name %q", id.ApplicationName),
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gallery Application", segmentsStr)
}

func (id GalleryApplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/applications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ApplicationName)
}

// GalleryApplicationID parses a GalleryApplication ID into an GalleryApplicationId struct
func GalleryApplicationID(input string) (*GalleryApplicationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := GalleryApplicationId{
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
	if resourceId.ApplicationName, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
