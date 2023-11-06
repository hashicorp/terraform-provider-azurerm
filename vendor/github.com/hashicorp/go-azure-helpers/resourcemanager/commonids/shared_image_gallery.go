// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SharedImageGalleryId{}

// SharedImageGalleryId is a struct representing the Resource ID for a Shared Image Gallery
type SharedImageGalleryId struct {
	SubscriptionId    string
	ResourceGroupName string
	GalleryName       string
}

// NewSharedImageGalleryID returns a new sharedImageGalleryId struct
func NewSharedImageGalleryID(subscriptionId string, resourceGroupName string, galleryName string) SharedImageGalleryId {
	return SharedImageGalleryId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GalleryName:       galleryName,
	}
}

// ParseSharedImageGalleryID parses 'input' into a sharedImageGalleryId
func ParseSharedImageGalleryID(input string) (*SharedImageGalleryId, error) {
	parser := resourceids.NewParserFromResourceIdType(SharedImageGalleryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SharedImageGalleryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "galleryName", *parsed)
	}

	return &id, nil
}

// ParseSharedImageGalleryIDInsensitively parses 'input' case-insensitively into a sharedImageGalleryId
// note: this method should only be used for API response data and not user input
func ParseSharedImageGalleryIDInsensitively(input string) (*SharedImageGalleryId, error) {
	parser := resourceids.NewParserFromResourceIdType(SharedImageGalleryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SharedImageGalleryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "galleryName", *parsed)
	}

	return &id, nil
}

// ValidateSharedImageGalleryID validates the ID of a Shared Image Gallery
func ValidateSharedImageGalleryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected string, got %T", input))
		return
	}

	if _, err := ParseSharedImageGalleryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Shared Image Gallery ID
func (id SharedImageGalleryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GalleryName)
}

// String returns a human-readable description of the Shared Image Gallery ID
func (id SharedImageGalleryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
	}
	return fmt.Sprintf("Shared Image Gallery (%s)", strings.Join(components, "\n"))
}

// Segments returns a slice of Resource ID Segments which comprise this Shared Image Gallery ID
func (id SharedImageGalleryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticGalleries", "galleries", "galleries"),
		resourceids.UserSpecifiedSegment("galleryName", "galleryValue"),
	}
}
