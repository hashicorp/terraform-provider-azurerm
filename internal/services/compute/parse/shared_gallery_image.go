// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SharedGalleryImageId{}

type SharedGalleryImageId struct {
	GalleryName string
	ImageName   string
}

func NewSharedGalleryImageID(galleryName, imageName string) SharedGalleryImageId {
	return SharedGalleryImageId{
		GalleryName: galleryName,
		ImageName:   imageName,
	}
}

func (id SharedGalleryImageId) String() string {
	segments := []string{
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Image Name %q", id.ImageName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Gallery", segmentsStr)
}

func (id SharedGalleryImageId) ID() string {
	fmtString := "/sharedGalleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.GalleryName, id.ImageName)
}

func (id SharedGalleryImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("sharedGalleries", "sharedGalleries", "sharedGalleries"),
		resourceids.UserSpecifiedSegment("galleryName", "myGalleryName"),
		resourceids.StaticSegment("images", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "myImageName"),
	}
}

// SharedGalleryImageID parses a SharedGalleryImage Unique ID into an SharedGalleryImageId struct
func SharedGalleryImageID(input string) (*SharedGalleryImageId, error) {
	id := SharedGalleryImageId{}
	parsed, err := resourceids.NewParserFromResourceIdType(id).Parse(input, false)
	if err != nil {
		return nil, err
	}

	var ok bool
	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "galleryName", *parsed)
	}
	if id.ImageName, ok = parsed.Parsed["imageName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "imageName", *parsed)
	}

	return &id, nil
}
