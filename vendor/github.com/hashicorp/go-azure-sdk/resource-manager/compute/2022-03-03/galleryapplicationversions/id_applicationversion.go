package galleryapplicationversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApplicationVersionId{}

// ApplicationVersionId is a struct representing the Resource ID for a Application Version
type ApplicationVersionId struct {
	SubscriptionId    string
	ResourceGroupName string
	GalleryName       string
	ApplicationName   string
	VersionName       string
}

// NewApplicationVersionID returns a new ApplicationVersionId struct
func NewApplicationVersionID(subscriptionId string, resourceGroupName string, galleryName string, applicationName string, versionName string) ApplicationVersionId {
	return ApplicationVersionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GalleryName:       galleryName,
		ApplicationName:   applicationName,
		VersionName:       versionName,
	}
}

// ParseApplicationVersionID parses 'input' into a ApplicationVersionId
func ParseApplicationVersionID(input string) (*ApplicationVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "galleryName", *parsed)
	}

	if id.ApplicationName, ok = parsed.Parsed["applicationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ParseApplicationVersionIDInsensitively parses 'input' case-insensitively into a ApplicationVersionId
// note: this method should only be used for API response data and not user input
func ParseApplicationVersionIDInsensitively(input string) (*ApplicationVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GalleryName, ok = parsed.Parsed["galleryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "galleryName", *parsed)
	}

	if id.ApplicationName, ok = parsed.Parsed["applicationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ValidateApplicationVersionID checks that 'input' can be parsed as a Application Version ID
func ValidateApplicationVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Version ID
func (id ApplicationVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/%s/applications/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GalleryName, id.ApplicationName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Version ID
func (id ApplicationVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticGalleries", "galleries", "galleries"),
		resourceids.UserSpecifiedSegment("galleryName", "galleryValue"),
		resourceids.StaticSegment("staticApplications", "applications", "applications"),
		resourceids.UserSpecifiedSegment("applicationName", "applicationValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

// String returns a human-readable description of this Application Version ID
func (id ApplicationVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gallery Name: %q", id.GalleryName),
		fmt.Sprintf("Application Name: %q", id.ApplicationName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Application Version (%s)", strings.Join(components, "\n"))
}
