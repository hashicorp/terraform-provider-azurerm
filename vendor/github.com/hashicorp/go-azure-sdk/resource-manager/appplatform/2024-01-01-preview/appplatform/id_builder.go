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
	recaser.RegisterResourceId(&BuilderId{})
}

var _ resourceids.ResourceId = &BuilderId{}

// BuilderId is a struct representing the Resource ID for a Builder
type BuilderId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	BuildServiceName  string
	BuilderName       string
}

// NewBuilderID returns a new BuilderId struct
func NewBuilderID(subscriptionId string, resourceGroupName string, springName string, buildServiceName string, builderName string) BuilderId {
	return BuilderId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		BuildServiceName:  buildServiceName,
		BuilderName:       builderName,
	}
}

// ParseBuilderID parses 'input' into a BuilderId
func ParseBuilderID(input string) (*BuilderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuilderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuilderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuilderIDInsensitively parses 'input' case-insensitively into a BuilderId
// note: this method should only be used for API response data and not user input
func ParseBuilderIDInsensitively(input string) (*BuilderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuilderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuilderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuilderId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BuilderName, ok = input.Parsed["builderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "builderName", input)
	}

	return nil
}

// ValidateBuilderID checks that 'input' can be parsed as a Builder ID
func ValidateBuilderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuilderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Builder ID
func (id BuilderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Builder ID
func (id BuilderId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBuilders", "builders", "builders"),
		resourceids.UserSpecifiedSegment("builderName", "builderName"),
	}
}

// String returns a human-readable description of this Builder ID
func (id BuilderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Build Service Name: %q", id.BuildServiceName),
		fmt.Sprintf("Builder Name: %q", id.BuilderName),
	}
	return fmt.Sprintf("Builder (%s)", strings.Join(components, "\n"))
}
