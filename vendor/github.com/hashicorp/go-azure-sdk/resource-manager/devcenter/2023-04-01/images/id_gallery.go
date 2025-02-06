package images

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GalleryId{})
}

var _ resourceids.ResourceId = &GalleryId{}

// GalleryId is a struct representing the Resource ID for a Gallery
type GalleryId struct {
	SubscriptionId    string
	ResourceGroupName string
	DevCenterName     string
	GalleryName       string
}

// NewGalleryID returns a new GalleryId struct
func NewGalleryID(subscriptionId string, resourceGroupName string, devCenterName string, galleryName string) GalleryId {
	return GalleryId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DevCenterName:     devCenterName,
		GalleryName:       galleryName,
	}
}

// ParseGalleryID parses 'input' into a GalleryId
func ParseGalleryID(input string) (*GalleryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GalleryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GalleryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGalleryIDInsensitively parses 'input' case-insensitively into a GalleryId
// note: this method should only be used for API response data and not user input
func ParseGalleryIDInsensitively(input string) (*GalleryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GalleryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GalleryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GalleryId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateGalleryID checks that 'input' can be parsed as a Gallery ID
func ValidateGalleryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGalleryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gallery ID
func (id GalleryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/galleries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.GalleryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gallery ID
func (id GalleryId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Gallery ID
func (id GalleryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
	}
	return fmt.Sprintf("Gallery (%s)", strings.Join(components, "\n"))
}
