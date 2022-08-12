package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

// SharedGalleryImageID parses a SharedGalleryImage Unique ID into an SharedGalleryImageId struct
func SharedGalleryImageID(input string) (*SharedGalleryImageId, error) {
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

	newParser := resourceids.NewParser(segments)

	id, err := newParser.Parse(input, false)
	if err != nil {
		return nil, err
	}

	resourceId := SharedGalleryImageId{
		GalleryName: id.Parsed["galleryName"],
		ImageName:   id.Parsed["imageName"],
	}

	if resourceId.GalleryName == "" {
		return nil, fmt.Errorf("ID was missing the 'GalleryName' element")
	}

	if resourceId.ImageName == "" {
		return nil, fmt.Errorf("ID was missing the 'ImageName' element")
	}

	return &resourceId, nil
}
