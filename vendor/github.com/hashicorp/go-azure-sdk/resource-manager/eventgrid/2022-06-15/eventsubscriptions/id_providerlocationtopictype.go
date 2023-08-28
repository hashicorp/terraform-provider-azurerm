package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderLocationTopicTypeId{}

// ProviderLocationTopicTypeId is a struct representing the Resource ID for a Provider Location Topic Type
type ProviderLocationTopicTypeId struct {
	SubscriptionId    string
	ResourceGroupName string
	LocationName      string
	TopicTypeName     string
}

// NewProviderLocationTopicTypeID returns a new ProviderLocationTopicTypeId struct
func NewProviderLocationTopicTypeID(subscriptionId string, resourceGroupName string, locationName string, topicTypeName string) ProviderLocationTopicTypeId {
	return ProviderLocationTopicTypeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LocationName:      locationName,
		TopicTypeName:     topicTypeName,
	}
}

// ParseProviderLocationTopicTypeID parses 'input' into a ProviderLocationTopicTypeId
func ParseProviderLocationTopicTypeID(input string) (*ProviderLocationTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLocationTopicTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLocationTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ParseProviderLocationTopicTypeIDInsensitively parses 'input' case-insensitively into a ProviderLocationTopicTypeId
// note: this method should only be used for API response data and not user input
func ParseProviderLocationTopicTypeIDInsensitively(input string) (*ProviderLocationTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLocationTopicTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLocationTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ValidateProviderLocationTopicTypeID checks that 'input' can be parsed as a Provider Location Topic Type ID
func ValidateProviderLocationTopicTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderLocationTopicTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Location Topic Type ID
func (id ProviderLocationTopicTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/locations/%s/topicTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.TopicTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Location Topic Type ID
func (id ProviderLocationTopicTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticTopicTypes", "topicTypes", "topicTypes"),
		resourceids.UserSpecifiedSegment("topicTypeName", "topicTypeValue"),
	}
}

// String returns a human-readable description of this Provider Location Topic Type ID
func (id ProviderLocationTopicTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Topic Type Name: %q", id.TopicTypeName),
	}
	return fmt.Sprintf("Provider Location Topic Type (%s)", strings.Join(components, "\n"))
}
