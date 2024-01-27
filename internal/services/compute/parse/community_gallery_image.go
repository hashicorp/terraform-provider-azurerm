// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &CommunityGalleryImageId{}

type CommunityGalleryImageId struct {
	GalleryName string
	ImageName   string
}

func NewCommunityGalleryImageID(galleryName, imageName string) CommunityGalleryImageId {
	return CommunityGalleryImageId{
		GalleryName: galleryName,
		ImageName:   imageName,
	}
}

func (id CommunityGalleryImageId) String() string {
	segments := []string{
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Image Name %q", id.ImageName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Community Gallery", segmentsStr)
}

func (id CommunityGalleryImageId) ID() string {
	fmtString := "/communityGalleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.GalleryName, id.ImageName)
}

func (id CommunityGalleryImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("communityGalleries", "communityGalleries", "communityGalleries"),
		resourceids.UserSpecifiedSegment("galleryName", "myGalleryName"),
		resourceids.StaticSegment("images", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "myImageName"),
	}
}

// CommunityGalleryImageID parses a CommunityGalleryImage Unique ID into an CommunityGalleryImageId struct
func CommunityGalleryImageID(input string) (*CommunityGalleryImageId, error) {
	id := CommunityGalleryImageId{}
	parsed, err := resourceids.NewParserFromResourceIdType(&id).Parse(input, false)
	if err != nil {
		return nil, err
	}

	var ok bool
	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(&id, "galleryName", *parsed)
	}
	if id.ImageName, ok = parsed.Parsed["imageName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(&id, "imageName", *parsed)
	}

	return &id, nil
}

func (id *CommunityGalleryImageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.GalleryName, ok = input.Parsed["galleryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "galleryName", input)
	}

	if id.ImageName, ok = input.Parsed["imageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageName", input)
	}

	return nil
}
