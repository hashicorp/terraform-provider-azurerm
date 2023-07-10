// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SharedGalleryImageVersionId{}

type SharedGalleryImageVersionId struct {
	GalleryName string
	ImageName   string
	Version     string
}

func NewSharedGalleryImageVersionID(galleryName, imageName, version string) SharedGalleryImageVersionId {
	return SharedGalleryImageVersionId{
		GalleryName: galleryName,
		ImageName:   imageName,
		Version:     version,
	}
}

func (id SharedGalleryImageVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Gallery Name %q", id.GalleryName),
		fmt.Sprintf("Image Name %q", id.ImageName),
		fmt.Sprintf("Version %q", id.Version),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Gallery", segmentsStr)
}

func (id SharedGalleryImageVersionId) ID() string {
	fmtString := "/sharedGalleries/%s/images/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.GalleryName, id.ImageName, id.Version)
}

func (id SharedGalleryImageVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("sharedGalleries", "sharedGalleries", "sharedGalleries"),
		resourceids.UserSpecifiedSegment("galleryName", "myGalleryName"),
		resourceids.StaticSegment("images", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "myImageName"),
		resourceids.StaticSegment("versions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("version", "myImageVersion"),
	}
}

// SharedGalleryImageVersionID parses a SharedGalleryImageVersion Unique ID into an SharedGalleryImageVersionId struct
func SharedGalleryImageVersionID(input string) (*SharedGalleryImageVersionId, error) {
	id := SharedGalleryImageVersionId{}
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
	if id.Version, ok = parsed.Parsed["version"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "version", *parsed)
	}

	// Additional validation for version, it can be the word "latest" or
	// a string in the format of Major.Minor.Patch, it must always be
	// a semantic version...

	if !strings.EqualFold(id.Version, "latest") {
		versionParts := strings.Split(id.Version, ".")

		if len(versionParts) != 3 {
			return nil, fmt.Errorf("ID 'Version' element is invalid, 'Version' must either be 'latest' or the semantic version(Major.Minor.Patch) for the image, got %s", id.Version)
		}

		for _, v := range versionParts {
			if _, err := strconv.Atoi(v); err != nil {
				return nil, fmt.Errorf("ID 'Version' element is invalid, semantic version elements must all be valid integers, got %s", id.Version)
			}
		}
	}

	return &id, nil
}
