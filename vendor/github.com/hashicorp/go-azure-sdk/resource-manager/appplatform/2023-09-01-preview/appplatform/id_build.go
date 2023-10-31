package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BuildId{}

// BuildId is a struct representing the Resource ID for a Build
type BuildId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	BuildServiceName  string
	BuildName         string
}

// NewBuildID returns a new BuildId struct
func NewBuildID(subscriptionId string, resourceGroupName string, springName string, buildServiceName string, buildName string) BuildId {
	return BuildId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		BuildServiceName:  buildServiceName,
		BuildName:         buildName,
	}
}

// ParseBuildID parses 'input' into a BuildId
func ParseBuildID(input string) (*BuildId, error) {
	parser := resourceids.NewParserFromResourceIdType(BuildId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BuildId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.BuildServiceName, ok = parsed.Parsed["buildServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildServiceName", *parsed)
	}

	if id.BuildName, ok = parsed.Parsed["buildName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildName", *parsed)
	}

	return &id, nil
}

// ParseBuildIDInsensitively parses 'input' case-insensitively into a BuildId
// note: this method should only be used for API response data and not user input
func ParseBuildIDInsensitively(input string) (*BuildId, error) {
	parser := resourceids.NewParserFromResourceIdType(BuildId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BuildId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.BuildServiceName, ok = parsed.Parsed["buildServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildServiceName", *parsed)
	}

	if id.BuildName, ok = parsed.Parsed["buildName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildName", *parsed)
	}

	return &id, nil
}

// ValidateBuildID checks that 'input' can be parsed as a Build ID
func ValidateBuildID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build ID
func (id BuildId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builds/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuildName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build ID
func (id BuildId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticBuildServices", "buildServices", "buildServices"),
		resourceids.UserSpecifiedSegment("buildServiceName", "buildServiceValue"),
		resourceids.StaticSegment("staticBuilds", "builds", "builds"),
		resourceids.UserSpecifiedSegment("buildName", "buildValue"),
	}
}

// String returns a human-readable description of this Build ID
func (id BuildId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Build Service Name: %q", id.BuildServiceName),
		fmt.Sprintf("Build Name: %q", id.BuildName),
	}
	return fmt.Sprintf("Build (%s)", strings.Join(components, "\n"))
}
