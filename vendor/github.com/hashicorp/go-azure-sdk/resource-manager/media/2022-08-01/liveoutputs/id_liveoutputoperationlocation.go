package liveoutputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LiveOutputOperationLocationId{}

// LiveOutputOperationLocationId is a struct representing the Resource ID for a Live Output Operation Location
type LiveOutputOperationLocationId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	LiveEventName     string
	LiveOutputName    string
	OperationId       string
}

// NewLiveOutputOperationLocationID returns a new LiveOutputOperationLocationId struct
func NewLiveOutputOperationLocationID(subscriptionId string, resourceGroupName string, mediaServiceName string, liveEventName string, liveOutputName string, operationId string) LiveOutputOperationLocationId {
	return LiveOutputOperationLocationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		LiveEventName:     liveEventName,
		LiveOutputName:    liveOutputName,
		OperationId:       operationId,
	}
}

// ParseLiveOutputOperationLocationID parses 'input' into a LiveOutputOperationLocationId
func ParseLiveOutputOperationLocationID(input string) (*LiveOutputOperationLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveOutputOperationLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LiveOutputOperationLocationId{}

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

	if id.LiveOutputName, ok = parsed.Parsed["liveOutputName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "liveOutputName", *parsed)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ParseLiveOutputOperationLocationIDInsensitively parses 'input' case-insensitively into a LiveOutputOperationLocationId
// note: this method should only be used for API response data and not user input
func ParseLiveOutputOperationLocationIDInsensitively(input string) (*LiveOutputOperationLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveOutputOperationLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LiveOutputOperationLocationId{}

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

	if id.LiveOutputName, ok = parsed.Parsed["liveOutputName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "liveOutputName", *parsed)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ValidateLiveOutputOperationLocationID checks that 'input' can be parsed as a Live Output Operation Location ID
func ValidateLiveOutputOperationLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLiveOutputOperationLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Live Output Operation Location ID
func (id LiveOutputOperationLocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/liveEvents/%s/liveOutputs/%s/operationLocations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.LiveEventName, id.LiveOutputName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Live Output Operation Location ID
func (id LiveOutputOperationLocationId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticLiveOutputs", "liveOutputs", "liveOutputs"),
		resourceids.UserSpecifiedSegment("liveOutputName", "liveOutputValue"),
		resourceids.StaticSegment("staticOperationLocations", "operationLocations", "operationLocations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Live Output Operation Location ID
func (id LiveOutputOperationLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Live Event Name: %q", id.LiveEventName),
		fmt.Sprintf("Live Output Name: %q", id.LiveOutputName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Live Output Operation Location (%s)", strings.Join(components, "\n"))
}
