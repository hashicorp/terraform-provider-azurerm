package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SpringId{}

// SpringId is a struct representing the Resource ID for a Spring
type SpringId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
}

// NewSpringID returns a new SpringId struct
func NewSpringID(subscriptionId string, resourceGroupName string, springName string) SpringId {
	return SpringId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
	}
}

// ParseSpringID parses 'input' into a SpringId
func ParseSpringID(input string) (*SpringId, error) {
	parser := resourceids.NewParserFromResourceIdType(SpringId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpringId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSpringIDInsensitively parses 'input' case-insensitively into a SpringId
// note: this method should only be used for API response data and not user input
func ParseSpringIDInsensitively(input string) (*SpringId, error) {
	parser := resourceids.NewParserFromResourceIdType(SpringId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpringId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SpringId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateSpringID checks that 'input' can be parsed as a Spring ID
func ValidateSpringID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSpringID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Spring ID
func (id SpringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName)
}

// Segments returns a slice of Resource ID Segments which comprise this Spring ID
func (id SpringId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
	}
}

// String returns a human-readable description of this Spring ID
func (id SpringId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
	}
	return fmt.Sprintf("Spring (%s)", strings.Join(components, "\n"))
}
