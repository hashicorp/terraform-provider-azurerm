package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

// SharedGalleryImageVersionID parses a SharedGalleryImageVersion Unique ID into an SharedGalleryImageVersionId struct
func SharedGalleryImageVersionID(input string) (*SharedGalleryImageVersionId, error) {
	segments := make([]resourceids.Segment, 0)

	segments = append(segments, resourceids.Segment{
		FixedValue: utils.String("sharedGalleries"),
		Name:       "sharedGalleries",
		Type:       resourceids.StaticSegmentType,
	})

	segments = append(segments, resourceids.Segment{
		ExampleValue: "myGalleryName",
		Name:         "galleryName",
		Type:         resourceids.UserSpecifiedSegmentType,
	})

	segments = append(segments, resourceids.Segment{
		FixedValue: utils.String("images"),
		Name:       "images",
		Type:       resourceids.StaticSegmentType,
	})

	segments = append(segments, resourceids.Segment{
		ExampleValue: "myImageName",
		Name:         "imageName",
		Type:         resourceids.UserSpecifiedSegmentType,
	})

	segments = append(segments, resourceids.Segment{
		FixedValue: utils.String("versions"),
		Name:       "versions",
		Type:       resourceids.StaticSegmentType,
	})

	segments = append(segments, resourceids.Segment{
		ExampleValue: "myImageVersion",
		Name:         "version",
		Type:         resourceids.UserSpecifiedSegmentType,
	})

	newParser := resourceids.NewParser(segments)

	id, err := newParser.Parse(input, false)
	if err != nil {
		return nil, err
	}

	resourceId := SharedGalleryImageVersionId{
		GalleryName: id.Parsed["galleryName"],
		ImageName:   id.Parsed["imageName"],
		Version:     id.Parsed["version"],
	}

	if resourceId.GalleryName == "" {
		return nil, fmt.Errorf("ID was missing the 'GalleryName' element")
	}

	if resourceId.ImageName == "" {
		return nil, fmt.Errorf("ID was missing the 'ImageName' element")
	}

	if resourceId.Version == "" {
		return nil, fmt.Errorf("ID was missing the 'Version' element")
	}

	// Additional validation for version, it can be the word "latest" or
	// a string in the format of Major.Minor.Patch, it must always be
	// a semantic version...

	if !strings.EqualFold(resourceId.Version, "latest") {
		versionParts := strings.Split(resourceId.Version, ".")

		if len(versionParts) != 3 {
			return nil, fmt.Errorf("ID 'Version' element is invalid, 'Version' must either be 'latest' or the semantic version(Major.Minor.Patch) for the image, got %s", resourceId.Version)
		}

		for _, v := range versionParts {
			if _, err := strconv.Atoi(v); err != nil {
				return nil, fmt.Errorf("ID 'Version' element is invalid, semantic version elements must all be valid integers, got %s", resourceId.Version)
			}
		}
	}

	return &resourceId, nil
}
