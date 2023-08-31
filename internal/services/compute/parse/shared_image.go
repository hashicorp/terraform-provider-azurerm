// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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
		fmt.Sprintf("Image Name %q", id.ImageName),
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Image", segmentsStr)
}

func (id SharedImageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ImageName)
}

// SharedImageID parses a SharedImage ID into an SharedImageId struct
func SharedImageID(input string) (*SharedImageId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SharedImage ID: %+v", input, err)
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
