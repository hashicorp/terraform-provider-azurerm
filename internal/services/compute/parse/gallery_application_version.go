package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GalleryApplicationVersionId struct {
	SubscriptionId  string
	ResourceGroup   string
	GalleryName     string
	ApplicationName string
	VersionName     string
}

func NewGalleryApplicationVersionID(subscriptionId, resourceGroup, galleryName, applicationName, versionName string) GalleryApplicationVersionId {
	return GalleryApplicationVersionId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		GalleryName:     galleryName,
		ApplicationName: applicationName,
		VersionName:     versionName,
	}
}

func (id GalleryApplicationVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Application Name %q", id.ApplicationName),
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gallery Application Version", segmentsStr)
}

func (id GalleryApplicationVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/applications/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName)
}

// GalleryApplicationVersionID parses a GalleryApplicationVersion ID into an GalleryApplicationVersionId struct
func GalleryApplicationVersionID(input string) (*GalleryApplicationVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := GalleryApplicationVersionId{
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
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
