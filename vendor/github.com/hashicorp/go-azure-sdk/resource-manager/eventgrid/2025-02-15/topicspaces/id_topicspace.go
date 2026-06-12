package topicspaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TopicSpaceId{})
}

var _ resourceids.ResourceId = &TopicSpaceId{}

// TopicSpaceId is a struct representing the Resource ID for a Topic Space
type TopicSpaceId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	TopicSpaceName    string
}

// NewTopicSpaceID returns a new TopicSpaceId struct
func NewTopicSpaceID(subscriptionId string, resourceGroupName string, namespaceName string, topicSpaceName string) TopicSpaceId {
	return TopicSpaceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		TopicSpaceName:    topicSpaceName,
	}
}

// ParseTopicSpaceID parses 'input' into a TopicSpaceId
func ParseTopicSpaceID(input string) (*TopicSpaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TopicSpaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TopicSpaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTopicSpaceIDInsensitively parses 'input' case-insensitively into a TopicSpaceId
// note: this method should only be used for API response data and not user input
func ParseTopicSpaceIDInsensitively(input string) (*TopicSpaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TopicSpaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TopicSpaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TopicSpaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.TopicSpaceName, ok = input.Parsed["topicSpaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "topicSpaceName", input)
	}

	return nil
}

// ValidateTopicSpaceID checks that 'input' can be parsed as a Topic Space ID
func ValidateTopicSpaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTopicSpaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Topic Space ID
func (id TopicSpaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/topicSpaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicSpaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Topic Space ID
func (id TopicSpaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticTopicSpaces", "topicSpaces", "topicSpaces"),
		resourceids.UserSpecifiedSegment("topicSpaceName", "topicSpaceName"),
	}
}

// String returns a human-readable description of this Topic Space ID
func (id TopicSpaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Space Name: %q", id.TopicSpaceName),
	}
	return fmt.Sprintf("Topic Space (%s)", strings.Join(components, "\n"))
}
