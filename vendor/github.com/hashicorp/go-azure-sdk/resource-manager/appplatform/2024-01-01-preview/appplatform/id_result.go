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
	recaser.RegisterResourceId(&ResultId{})
}

var _ resourceids.ResourceId = &ResultId{}

// ResultId is a struct representing the Resource ID for a Result
type ResultId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	BuildServiceName  string
	BuildName         string
	ResultName        string
}

// NewResultID returns a new ResultId struct
func NewResultID(subscriptionId string, resourceGroupName string, springName string, buildServiceName string, buildName string, resultName string) ResultId {
	return ResultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		BuildServiceName:  buildServiceName,
		BuildName:         buildName,
		ResultName:        resultName,
	}
}

// ParseResultID parses 'input' into a ResultId
func ParseResultID(input string) (*ResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResultIDInsensitively parses 'input' case-insensitively into a ResultId
// note: this method should only be used for API response data and not user input
func ParseResultIDInsensitively(input string) (*ResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResultId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BuildName, ok = input.Parsed["buildName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "buildName", input)
	}

	if id.ResultName, ok = input.Parsed["resultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resultName", input)
	}

	return nil
}

// ValidateResultID checks that 'input' can be parsed as a Result ID
func ValidateResultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Result ID
func (id ResultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builds/%s/results/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuildName, id.ResultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Result ID
func (id ResultId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBuilds", "builds", "builds"),
		resourceids.UserSpecifiedSegment("buildName", "buildName"),
		resourceids.StaticSegment("staticResults", "results", "results"),
		resourceids.UserSpecifiedSegment("resultName", "resultName"),
	}
}

// String returns a human-readable description of this Result ID
func (id ResultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Build Service Name: %q", id.BuildServiceName),
		fmt.Sprintf("Build Name: %q", id.BuildName),
		fmt.Sprintf("Result Name: %q", id.ResultName),
	}
	return fmt.Sprintf("Result (%s)", strings.Join(components, "\n"))
}
