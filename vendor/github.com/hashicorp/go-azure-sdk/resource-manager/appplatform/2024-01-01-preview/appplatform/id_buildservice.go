package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BuildServiceId{})
}

var _ resourceids.ResourceId = &BuildServiceId{}

// BuildServiceId is a struct representing the Resource ID for a Build Service
type BuildServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	BuildServiceName  string
}

// NewBuildServiceID returns a new BuildServiceId struct
func NewBuildServiceID(subscriptionId string, resourceGroupName string, springName string, buildServiceName string) BuildServiceId {
	return BuildServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		BuildServiceName:  buildServiceName,
	}
}

// ParseBuildServiceID parses 'input' into a BuildServiceId
func ParseBuildServiceID(input string) (*BuildServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuildServiceIDInsensitively parses 'input' case-insensitively into a BuildServiceId
// note: this method should only be used for API response data and not user input
func ParseBuildServiceIDInsensitively(input string) (*BuildServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuildServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.BuildServiceName, ok = input.Parsed["buildServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "buildServiceName", input)
	}

	return nil
}

// ValidateBuildServiceID checks that 'input' can be parsed as a Build Service ID
func ValidateBuildServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build Service ID
func (id BuildServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build Service ID
func (id BuildServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticBuildServices", "buildServices", "buildServices"),
		resourceids.UserSpecifiedSegment("buildServiceName", "buildServiceName"),
	}
}

// String returns a human-readable description of this Build Service ID
func (id BuildServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Build Service Name: %q", id.BuildServiceName),
	}
	return fmt.Sprintf("Build Service (%s)", strings.Join(components, "\n"))
}
