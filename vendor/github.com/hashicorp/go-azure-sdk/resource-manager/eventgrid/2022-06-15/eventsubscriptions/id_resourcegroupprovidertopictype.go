package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ResourceGroupProviderTopicTypeId{}

// ResourceGroupProviderTopicTypeId is a struct representing the Resource ID for a Resource Group Provider Topic Type
type ResourceGroupProviderTopicTypeId struct {
	SubscriptionId    string
	ResourceGroupName string
	TopicTypeName     string
}

// NewResourceGroupProviderTopicTypeID returns a new ResourceGroupProviderTopicTypeId struct
func NewResourceGroupProviderTopicTypeID(subscriptionId string, resourceGroupName string, topicTypeName string) ResourceGroupProviderTopicTypeId {
	return ResourceGroupProviderTopicTypeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		TopicTypeName:     topicTypeName,
	}
}

// ParseResourceGroupProviderTopicTypeID parses 'input' into a ResourceGroupProviderTopicTypeId
func ParseResourceGroupProviderTopicTypeID(input string) (*ResourceGroupProviderTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGroupProviderTopicTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGroupProviderTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ParseResourceGroupProviderTopicTypeIDInsensitively parses 'input' case-insensitively into a ResourceGroupProviderTopicTypeId
// note: this method should only be used for API response data and not user input
func ParseResourceGroupProviderTopicTypeIDInsensitively(input string) (*ResourceGroupProviderTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGroupProviderTopicTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGroupProviderTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ValidateResourceGroupProviderTopicTypeID checks that 'input' can be parsed as a Resource Group Provider Topic Type ID
func ValidateResourceGroupProviderTopicTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGroupProviderTopicTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Group Provider Topic Type ID
func (id ResourceGroupProviderTopicTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/topicTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TopicTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Group Provider Topic Type ID
func (id ResourceGroupProviderTopicTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticTopicTypes", "topicTypes", "topicTypes"),
		resourceids.UserSpecifiedSegment("topicTypeName", "topicTypeValue"),
	}
}

// String returns a human-readable description of this Resource Group Provider Topic Type ID
func (id ResourceGroupProviderTopicTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Topic Type Name: %q", id.TopicTypeName),
	}
	return fmt.Sprintf("Resource Group Provider Topic Type (%s)", strings.Join(components, "\n"))
}
