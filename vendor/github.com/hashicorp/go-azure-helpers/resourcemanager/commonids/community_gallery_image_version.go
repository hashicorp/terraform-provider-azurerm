// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &CommunityGalleryImageVersionId{}

// CommunityGalleryImageVersionId is a struct representing the Resource ID for a Community Gallery Image Version
type CommunityGalleryImageVersionId struct {
	CommunityGalleryName string
	ImageName            string
	VersionName          string
}

// NewCommunityGalleryImageVersionID returns a new CommunityGalleryImageVersionId struct
func NewCommunityGalleryImageVersionID(communityGallery string, image string, version string) CommunityGalleryImageVersionId {
	return CommunityGalleryImageVersionId{
		CommunityGalleryName: communityGallery,
		ImageName:            image,
		VersionName:          version,
	}
}

// ParseCommunityGalleryImageVersionID parses 'input' into a CommunityGalleryImageVersionId
func ParseCommunityGalleryImageVersionID(input string) (*CommunityGalleryImageVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunityGalleryImageVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunityGalleryImageVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCommunityGalleryImageVersionIDInsensitively parses 'input' case-insensitively into a CommunityGalleryImageVersionId
// note: this method should only be used for API response data and not user input
func ParseCommunityGalleryImageVersionIDInsensitively(input string) (*CommunityGalleryImageVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunityGalleryImageVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunityGalleryImageVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CommunityGalleryImageVersionId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.CommunityGalleryName, ok = input.Parsed["communityGalleryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "communityGalleryName", input)
	}

	if id.ImageName, ok = input.Parsed["imageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageName", input)
	}

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// ValidateCommunityGalleryImageVersionID checks that 'input' can be parsed as a Community Gallery Image Version ID
func ValidateCommunityGalleryImageVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommunityGalleryImageVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Community Gallery Image Version ID
func (id CommunityGalleryImageVersionId) ID() string {
	fmtString := "/communityGalleries/%s/images/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.CommunityGalleryName, id.ImageName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Community Gallery Image Version ID
func (id CommunityGalleryImageVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticCommunityGalleries", "communityGalleries", "communityGalleries"),
		resourceids.UserSpecifiedSegment("communityGalleryName", "communityGalleryValue"),
		resourceids.StaticSegment("staticImages", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "imageValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

// String returns a human-readable description of this Community Gallery Image Version ID
func (id CommunityGalleryImageVersionId) String() string {
	components := []string{
		fmt.Sprintf("Community Gallery Name: %q", id.CommunityGalleryName),
		fmt.Sprintf("Image Name: %q", id.ImageName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Community Gallery Image Version (%s)", strings.Join(components, "\n"))
}
