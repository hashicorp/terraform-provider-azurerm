package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

// CommunityGalleryImageID parses a CommunityGalleryImage Unique ID into an CommunityGalleryImageId struct
func CommunityGalleryImageID(input string) (*CommunityGalleryImageId, error) {
	segments := make([]resourceids.Segment, 0)

	segments = append(segments, resourceids.Segment{
		FixedValue: utils.String("communityGalleries"),
		Name:       "communityGalleries",
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

	resourceId := CommunityGalleryImageId{
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
