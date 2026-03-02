package imageversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ImageVersionId{})
}

var _ resourceids.ResourceId = &ImageVersionId{}

// ImageVersionId is a struct representing the Resource ID for a Image Version
type ImageVersionId struct {
	SubscriptionId    string
	ResourceGroupName string
	DevCenterName     string
	GalleryName       string
	ImageName         string
	VersionName       string
}

// NewImageVersionID returns a new ImageVersionId struct
func NewImageVersionID(subscriptionId string, resourceGroupName string, devCenterName string, galleryName string, imageName string, versionName string) ImageVersionId {
	return ImageVersionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DevCenterName:     devCenterName,
		GalleryName:       galleryName,
		ImageName:         imageName,
		VersionName:       versionName,
	}
}

// ParseImageVersionID parses 'input' into a ImageVersionId
func ParseImageVersionID(input string) (*ImageVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseImageVersionIDInsensitively parses 'input' case-insensitively into a ImageVersionId
// note: this method should only be used for API response data and not user input
func ParseImageVersionIDInsensitively(input string) (*ImageVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ImageVersionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// ValidateImageVersionID checks that 'input' can be parsed as a Image Version ID
func ValidateImageVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImageVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Image Version ID
func (id ImageVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/galleries/%s/images/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.GalleryName, id.ImageName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Image Version ID
func (id ImageVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticGalleries", "galleries", "galleries"),
		resourceids.UserSpecifiedSegment("galleryName", "galleryName"),
		resourceids.StaticSegment("staticImages", "images", "images"),
		resourceids.UserSpecifiedSegment("imageName", "imageName"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
	}
}

// String returns a human-readable description of this Image Version ID
func (id ImageVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
		fmt.Sprintf("Image Name: %q", id.ImageName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Image Version (%s)", strings.Join(components, "\n"))
}
