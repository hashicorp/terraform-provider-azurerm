package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BuildPackBindingId{}

// BuildPackBindingId is a struct representing the Resource ID for a Build Pack Binding
type BuildPackBindingId struct {
	SubscriptionId       string
	ResourceGroupName    string
	SpringName           string
	BuildServiceName     string
	BuilderName          string
	BuildPackBindingName string
}

// NewBuildPackBindingID returns a new BuildPackBindingId struct
func NewBuildPackBindingID(subscriptionId string, resourceGroupName string, springName string, buildServiceName string, builderName string, buildPackBindingName string) BuildPackBindingId {
	return BuildPackBindingId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		SpringName:           springName,
		BuildServiceName:     buildServiceName,
		BuilderName:          builderName,
		BuildPackBindingName: buildPackBindingName,
	}
}

// ParseBuildPackBindingID parses 'input' into a BuildPackBindingId
func ParseBuildPackBindingID(input string) (*BuildPackBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(BuildPackBindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BuildPackBindingId{}

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

	if id.BuilderName, ok = parsed.Parsed["builderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "builderName", *parsed)
	}

	if id.BuildPackBindingName, ok = parsed.Parsed["buildPackBindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildPackBindingName", *parsed)
	}

	return &id, nil
}

// ParseBuildPackBindingIDInsensitively parses 'input' case-insensitively into a BuildPackBindingId
// note: this method should only be used for API response data and not user input
func ParseBuildPackBindingIDInsensitively(input string) (*BuildPackBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(BuildPackBindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BuildPackBindingId{}

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

	if id.BuilderName, ok = parsed.Parsed["builderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "builderName", *parsed)
	}

	if id.BuildPackBindingName, ok = parsed.Parsed["buildPackBindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "buildPackBindingName", *parsed)
	}

	return &id, nil
}

// ValidateBuildPackBindingID checks that 'input' can be parsed as a Build Pack Binding ID
func ValidateBuildPackBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildPackBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build Pack Binding ID
func (id BuildPackBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builders/%s/buildPackBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build Pack Binding ID
func (id BuildPackBindingId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBuilders", "builders", "builders"),
		resourceids.UserSpecifiedSegment("builderName", "builderValue"),
		resourceids.StaticSegment("staticBuildPackBindings", "buildPackBindings", "buildPackBindings"),
		resourceids.UserSpecifiedSegment("buildPackBindingName", "buildPackBindingValue"),
	}
}

// String returns a human-readable description of this Build Pack Binding ID
func (id BuildPackBindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Build Service Name: %q", id.BuildServiceName),
		fmt.Sprintf("Builder Name: %q", id.BuilderName),
		fmt.Sprintf("Build Pack Binding Name: %q", id.BuildPackBindingName),
	}
	return fmt.Sprintf("Build Pack Binding (%s)", strings.Join(components, "\n"))
}
