package liveoutputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LiveEventId{}

// LiveEventId is a struct representing the Resource ID for a Live Event
type LiveEventId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	LiveEventName     string
}

// NewLiveEventID returns a new LiveEventId struct
func NewLiveEventID(subscriptionId string, resourceGroupName string, mediaServiceName string, liveEventName string) LiveEventId {
	return LiveEventId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		LiveEventName:     liveEventName,
	}
}

// ParseLiveEventID parses 'input' into a LiveEventId
func ParseLiveEventID(input string) (*LiveEventId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveEventId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LiveEventId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.LiveEventName, ok = parsed.Parsed["liveEventName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "liveEventName", *parsed)
	}

	return &id, nil
}

// ParseLiveEventIDInsensitively parses 'input' case-insensitively into a LiveEventId
// note: this method should only be used for API response data and not user input
func ParseLiveEventIDInsensitively(input string) (*LiveEventId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveEventId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LiveEventId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.LiveEventName, ok = parsed.Parsed["liveEventName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "liveEventName", *parsed)
	}

	return &id, nil
}

// ValidateLiveEventID checks that 'input' can be parsed as a Live Event ID
func ValidateLiveEventID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLiveEventID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Live Event ID
func (id LiveEventId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/liveEvents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.LiveEventName)
}

// Segments returns a slice of Resource ID Segments which comprise this Live Event ID
func (id LiveEventId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticLiveEvents", "liveEvents", "liveEvents"),
		resourceids.UserSpecifiedSegment("liveEventName", "liveEventValue"),
	}
}

// String returns a human-readable description of this Live Event ID
func (id LiveEventId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Live Event Name: %q", id.LiveEventName),
	}
	return fmt.Sprintf("Live Event (%s)", strings.Join(components, "\n"))
}
