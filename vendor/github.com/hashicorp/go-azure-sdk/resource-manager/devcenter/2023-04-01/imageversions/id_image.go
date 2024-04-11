package imageversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ImageId{}

// ImageId is a struct representing the Resource ID for a Image
type ImageId struct {
	SubscriptionId    string
	ResourceGroupName string
	DevCenterName     string
	GalleryName       string
	ImageName         string
}

// NewImageID returns a new ImageId struct
func NewImageID(subscriptionId string, resourceGroupName string, devCenterName string, galleryName string, imageName string) ImageId {
	return ImageId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DevCenterName:     devCenterName,
		GalleryName:       galleryName,
		ImageName:         imageName,
	}
}

// ParseImageID parses 'input' into a ImageId
func ParseImageID(input string) (*ImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseImageIDInsensitively parses 'input' case-insensitively into a ImageId
// note: this method should only be used for API response data and not user input
func ParseImageIDInsensitively(input string) (*ImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ImageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.GalleryName, ok = input.Parsed["galleryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "galleryName", input)
	}

	if id.ImageName, ok = input.Parsed["imageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageName", input)
	}

	return nil
}

// ValidateImageID checks that 'input' can be parsed as a Image ID
func ValidateImageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Image ID
func (id ImageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/galleries/%s/images/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.GalleryName, id.ImageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Image ID
func (id ImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterValue"),
		resourceids.StaticSegment("staticGalleries", "galleries", "galleries"),
		resourceids.UserSpecifiedSegment("galleryName", "galleryValue"),
		resourceids.StaticSegment("staticImages", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "imageValue"),
	}
}

// String returns a human-readable description of this Image ID
func (id ImageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
		fmt.Sprintf("Image Name: %q", id.ImageName),
	}
	return fmt.Sprintf("Image (%s)", strings.Join(components, "\n"))
}
