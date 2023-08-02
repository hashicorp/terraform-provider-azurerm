package systemtopics

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SystemTopicId{}

// SystemTopicId is a struct representing the Resource ID for a System Topic
type SystemTopicId struct {
	SubscriptionId    string
	ResourceGroupName string
	SystemTopicName   string
}

// NewSystemTopicID returns a new SystemTopicId struct
func NewSystemTopicID(subscriptionId string, resourceGroupName string, systemTopicName string) SystemTopicId {
	return SystemTopicId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SystemTopicName:   systemTopicName,
	}
}

// ParseSystemTopicID parses 'input' into a SystemTopicId
func ParseSystemTopicID(input string) (*SystemTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(SystemTopicId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SystemTopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SystemTopicName, ok = parsed.Parsed["systemTopicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "systemTopicName", *parsed)
	}

	return &id, nil
}

// ParseSystemTopicIDInsensitively parses 'input' case-insensitively into a SystemTopicId
// note: this method should only be used for API response data and not user input
func ParseSystemTopicIDInsensitively(input string) (*SystemTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(SystemTopicId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SystemTopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SystemTopicName, ok = parsed.Parsed["systemTopicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "systemTopicName", *parsed)
	}

	return &id, nil
}

// ValidateSystemTopicID checks that 'input' can be parsed as a System Topic ID
func ValidateSystemTopicID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSystemTopicID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted System Topic ID
func (id SystemTopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/systemTopics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SystemTopicName)
}

// Segments returns a slice of Resource ID Segments which comprise this System Topic ID
func (id SystemTopicId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticSystemTopics", "systemTopics", "systemTopics"),
		resourceids.UserSpecifiedSegment("systemTopicName", "systemTopicValue"),
	}
}

// String returns a human-readable description of this System Topic ID
func (id SystemTopicId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("System Topic Name: %q", id.SystemTopicName),
	}
	return fmt.Sprintf("System Topic (%s)", strings.Join(components, "\n"))
}
