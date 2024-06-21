package galleryimageversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GalleryImageId{})
}

var _ resourceids.ResourceId = &GalleryImageId{}

// GalleryImageId is a struct representing the Resource ID for a Gallery Image
type GalleryImageId struct {
	SubscriptionId    string
	ResourceGroupName string
	GalleryName       string
	ImageName         string
}

// NewGalleryImageID returns a new GalleryImageId struct
func NewGalleryImageID(subscriptionId string, resourceGroupName string, galleryName string, imageName string) GalleryImageId {
	return GalleryImageId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GalleryName:       galleryName,
		ImageName:         imageName,
	}
}

// ParseGalleryImageID parses 'input' into a GalleryImageId
func ParseGalleryImageID(input string) (*GalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GalleryImageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GalleryImageId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGalleryImageIDInsensitively parses 'input' case-insensitively into a GalleryImageId
// note: this method should only be used for API response data and not user input
func ParseGalleryImageIDInsensitively(input string) (*GalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GalleryImageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GalleryImageId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GalleryImageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GalleryName, ok = input.Parsed["galleryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "galleryName", input)
	}

	if id.ImageName, ok = input.Parsed["imageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageName", input)
	}

	return nil
}

// ValidateGalleryImageID checks that 'input' can be parsed as a Gallery Image ID
func ValidateGalleryImageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGalleryImageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gallery Image ID
func (id GalleryImageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GalleryName, id.ImageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gallery Image ID
func (id GalleryImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticGalleries", "galleries", "galleries"),
		resourceids.UserSpecifiedSegment("galleryName", "galleryValue"),
		resourceids.StaticSegment("staticImages", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "imageValue"),
	}
}

// String returns a human-readable description of this Gallery Image ID
func (id GalleryImageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
		fmt.Sprintf("Image Name: %q", id.ImageName),
	}
	return fmt.Sprintf("Gallery Image (%s)", strings.Join(components, "\n"))
}
