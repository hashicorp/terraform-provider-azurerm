package galleryapplicationversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationVersionId{})
}

var _ resourceids.ResourceId = &ApplicationVersionId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ApplicationVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationVersionIDInsensitively parses 'input' case-insensitively into a ApplicationVersionId
// note: this method should only be used for API response data and not user input
func ParseApplicationVersionIDInsensitively(input string) (*ApplicationVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationVersionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApplicationName, ok = input.Parsed["applicationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationName", input)
	}

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("galleryName", "galleryName"),
		resourceids.StaticSegment("staticApplications", "applications", "applications"),
		resourceids.UserSpecifiedSegment("applicationName", "applicationName"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
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
