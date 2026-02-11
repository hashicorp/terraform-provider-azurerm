package namespacetopics

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NamespaceTopicId{})
}

var _ resourceids.ResourceId = &NamespaceTopicId{}

// NamespaceTopicId is a struct representing the Resource ID for a Namespace Topic
type NamespaceTopicId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	TopicName         string
}

// NewNamespaceTopicID returns a new NamespaceTopicId struct
func NewNamespaceTopicID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string) NamespaceTopicId {
	return NamespaceTopicId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		TopicName:         topicName,
	}
}

// ParseNamespaceTopicID parses 'input' into a NamespaceTopicId
func ParseNamespaceTopicID(input string) (*NamespaceTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceTopicId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceTopicId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNamespaceTopicIDInsensitively parses 'input' case-insensitively into a NamespaceTopicId
// note: this method should only be used for API response data and not user input
func ParseNamespaceTopicIDInsensitively(input string) (*NamespaceTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceTopicId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceTopicId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NamespaceTopicId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TopicName, ok = input.Parsed["topicName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "topicName", input)
	}

	return nil
}

// ValidateNamespaceTopicID checks that 'input' can be parsed as a Namespace Topic ID
func ValidateNamespaceTopicID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNamespaceTopicID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Namespace Topic ID
func (id NamespaceTopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/topics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName)
}

// Segments returns a slice of Resource ID Segments which comprise this Namespace Topic ID
func (id NamespaceTopicId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicName"),
	}
}

// String returns a human-readable description of this Namespace Topic ID
func (id NamespaceTopicId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
	}
	return fmt.Sprintf("Namespace Topic (%s)", strings.Join(components, "\n"))
}
