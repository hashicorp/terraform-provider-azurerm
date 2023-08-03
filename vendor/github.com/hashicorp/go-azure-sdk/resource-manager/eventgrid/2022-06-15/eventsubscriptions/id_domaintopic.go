package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DomainTopicId{}

// DomainTopicId is a struct representing the Resource ID for a Domain Topic
type DomainTopicId struct {
	SubscriptionId    string
	ResourceGroupName string
	DomainName        string
	TopicName         string
}

// NewDomainTopicID returns a new DomainTopicId struct
func NewDomainTopicID(subscriptionId string, resourceGroupName string, domainName string, topicName string) DomainTopicId {
	return DomainTopicId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DomainName:        domainName,
		TopicName:         topicName,
	}
}

// ParseDomainTopicID parses 'input' into a DomainTopicId
func ParseDomainTopicID(input string) (*DomainTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainTopicId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainTopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	return &id, nil
}

// ParseDomainTopicIDInsensitively parses 'input' case-insensitively into a DomainTopicId
// note: this method should only be used for API response data and not user input
func ParseDomainTopicIDInsensitively(input string) (*DomainTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainTopicId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainTopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	return &id, nil
}

// ValidateDomainTopicID checks that 'input' can be parsed as a Domain Topic ID
func ValidateDomainTopicID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDomainTopicID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Domain Topic ID
func (id DomainTopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/domains/%s/topics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DomainName, id.TopicName)
}

// Segments returns a slice of Resource ID Segments which comprise this Domain Topic ID
func (id DomainTopicId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainValue"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicValue"),
	}
}

// String returns a human-readable description of this Domain Topic ID
func (id DomainTopicId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
	}
	return fmt.Sprintf("Domain Topic (%s)", strings.Join(components, "\n"))
}
