// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &CommunityGalleryImageId{}

// CommunityGalleryImageId is a struct representing the Resource ID for a Community Gallery Image
type CommunityGalleryImageId struct {
	CommunityGalleryName string
	ImageName            string
}

// NewCommunityGalleryImageID returns a new CommunityGalleryImageId struct
func NewCommunityGalleryImageID(communityGallery string, image string) CommunityGalleryImageId {
	return CommunityGalleryImageId{
		CommunityGalleryName: communityGallery,
		ImageName:            image,
	}
}

// ParseCommunityGalleryImageID parses 'input' into a CommunityGalleryImageId
func ParseCommunityGalleryImageID(input string) (*CommunityGalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunityGalleryImageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunityGalleryImageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCommunityGalleryImageIDInsensitively parses 'input' case-insensitively into a CommunityGalleryImageId
// note: this method should only be used for API response data and not user input
func ParseCommunityGalleryImageIDInsensitively(input string) (*CommunityGalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunityGalleryImageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunityGalleryImageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CommunityGalleryImageId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.CommunityGalleryName, ok = input.Parsed["communityGalleryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "communityGalleryName", input)
	}

	if id.ImageName, ok = input.Parsed["imageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageName", input)
	}

	return nil
}

// ValidateCommunityGalleryImageID checks that 'input' can be parsed as a Community Gallery Image ID
func ValidateCommunityGalleryImageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommunityGalleryImageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Community Gallery Image ID
func (id CommunityGalleryImageId) ID() string {
	fmtString := "/communityGalleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.CommunityGalleryName, id.ImageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Community Gallery Image ID
func (id CommunityGalleryImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticCommunityGalleries", "communityGalleries", "communityGalleries"),
		resourceids.UserSpecifiedSegment("communityGalleryName", "communityGalleryValue"),
		resourceids.StaticSegment("staticImages", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "imageValue"),
	}
}

// String returns a human-readable description of this Community Gallery Image ID
func (id CommunityGalleryImageId) String() string {
	components := []string{
		fmt.Sprintf("Community Gallery Name: %q", id.CommunityGalleryName),
		fmt.Sprintf("Image Name: %q", id.ImageName),
	}
	return fmt.Sprintf("Community Gallery Image (%s)", strings.Join(components, "\n"))
}
