package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocationTopicTypeId{}

// LocationTopicTypeId is a struct representing the Resource ID for a Location Topic Type
type LocationTopicTypeId struct {
	SubscriptionId string
	LocationName   string
	TopicTypeName  string
}

// NewLocationTopicTypeID returns a new LocationTopicTypeId struct
func NewLocationTopicTypeID(subscriptionId string, locationName string, topicTypeName string) LocationTopicTypeId {
	return LocationTopicTypeId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		TopicTypeName:  topicTypeName,
	}
}

// ParseLocationTopicTypeID parses 'input' into a LocationTopicTypeId
func ParseLocationTopicTypeID(input string) (*LocationTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationTopicTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ParseLocationTopicTypeIDInsensitively parses 'input' case-insensitively into a LocationTopicTypeId
// note: this method should only be used for API response data and not user input
func ParseLocationTopicTypeIDInsensitively(input string) (*LocationTopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationTopicTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationTopicTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.TopicTypeName, ok = parsed.Parsed["topicTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", *parsed)
	}

	return &id, nil
}

// ValidateLocationTopicTypeID checks that 'input' can be parsed as a Location Topic Type ID
func ValidateLocationTopicTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationTopicTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location Topic Type ID
func (id LocationTopicTypeId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.EventGrid/locations/%s/topicTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.TopicTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Location Topic Type ID
func (id LocationTopicTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticTopicTypes", "topicTypes", "topicTypes"),
		resourceids.UserSpecifiedSegment("topicTypeName", "topicTypeValue"),
	}
}

// String returns a human-readable description of this Location Topic Type ID
func (id LocationTopicTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Topic Type Name: %q", id.TopicTypeName),
	}
	return fmt.Sprintf("Location Topic Type (%s)", strings.Join(components, "\n"))
}
