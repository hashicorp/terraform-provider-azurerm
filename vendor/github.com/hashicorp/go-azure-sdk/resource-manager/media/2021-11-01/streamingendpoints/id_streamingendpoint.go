package streamingendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StreamingEndpointId{}

// StreamingEndpointId is a struct representing the Resource ID for a Streaming Endpoint
type StreamingEndpointId struct {
	SubscriptionId        string
	ResourceGroupName     string
	MediaServiceName      string
	StreamingEndpointName string
}

// NewStreamingEndpointID returns a new StreamingEndpointId struct
func NewStreamingEndpointID(subscriptionId string, resourceGroupName string, mediaServiceName string, streamingEndpointName string) StreamingEndpointId {
	return StreamingEndpointId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		MediaServiceName:      mediaServiceName,
		StreamingEndpointName: streamingEndpointName,
	}
}

// ParseStreamingEndpointID parses 'input' into a StreamingEndpointId
func ParseStreamingEndpointID(input string) (*StreamingEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.StreamingEndpointName, ok = parsed.Parsed["streamingEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingEndpointName", *parsed)
	}

	return &id, nil
}

// ParseStreamingEndpointIDInsensitively parses 'input' case-insensitively into a StreamingEndpointId
// note: this method should only be used for API response data and not user input
func ParseStreamingEndpointIDInsensitively(input string) (*StreamingEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.StreamingEndpointName, ok = parsed.Parsed["streamingEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingEndpointName", *parsed)
	}

	return &id, nil
}

// ValidateStreamingEndpointID checks that 'input' can be parsed as a Streaming Endpoint ID
func ValidateStreamingEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStreamingEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Streaming Endpoint ID
func (id StreamingEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/streamingEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.StreamingEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Streaming Endpoint ID
func (id StreamingEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticStreamingEndpoints", "streamingEndpoints", "streamingEndpoints"),
		resourceids.UserSpecifiedSegment("streamingEndpointName", "streamingEndpointValue"),
	}
}

// String returns a human-readable description of this Streaming Endpoint ID
func (id StreamingEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Streaming Endpoint Name: %q", id.StreamingEndpointName),
	}
	return fmt.Sprintf("Streaming Endpoint (%s)", strings.Join(components, "\n"))
}
