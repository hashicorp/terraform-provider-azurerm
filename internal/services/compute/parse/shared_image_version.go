// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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

func (id SharedImageVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Image Name %q", id.ImageName),
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Image Version", segmentsStr)
}

func (id SharedImageVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName)
}

// SharedImageVersionID parses a SharedImageVersion ID into an SharedImageVersionId struct
func SharedImageVersionID(input string) (*SharedImageVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SharedImageVersion ID: %+v", input, err)
	}

	resourceId := SharedImageVersionId{
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
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
