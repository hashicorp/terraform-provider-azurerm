package streamingpoliciesandstreaminglocators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StreamingPolicyId{}

// StreamingPolicyId is a struct representing the Resource ID for a Streaming Policy
type StreamingPolicyId struct {
	SubscriptionId      string
	ResourceGroupName   string
	MediaServiceName    string
	StreamingPolicyName string
}

// NewStreamingPolicyID returns a new StreamingPolicyId struct
func NewStreamingPolicyID(subscriptionId string, resourceGroupName string, mediaServiceName string, streamingPolicyName string) StreamingPolicyId {
	return StreamingPolicyId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		MediaServiceName:    mediaServiceName,
		StreamingPolicyName: streamingPolicyName,
	}
}

// ParseStreamingPolicyID parses 'input' into a StreamingPolicyId
func ParseStreamingPolicyID(input string) (*StreamingPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.StreamingPolicyName, ok = parsed.Parsed["streamingPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingPolicyName", *parsed)
	}

	return &id, nil
}

// ParseStreamingPolicyIDInsensitively parses 'input' case-insensitively into a StreamingPolicyId
// note: this method should only be used for API response data and not user input
func ParseStreamingPolicyIDInsensitively(input string) (*StreamingPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.StreamingPolicyName, ok = parsed.Parsed["streamingPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateStreamingPolicyID checks that 'input' can be parsed as a Streaming Policy ID
func ValidateStreamingPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStreamingPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Streaming Policy ID
func (id StreamingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/streamingPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.StreamingPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Streaming Policy ID
func (id StreamingPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticStreamingPolicies", "streamingPolicies", "streamingPolicies"),
		resourceids.UserSpecifiedSegment("streamingPolicyName", "streamingPolicyValue"),
	}
}

// String returns a human-readable description of this Streaming Policy ID
func (id StreamingPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Streaming Policy Name: %q", id.StreamingPolicyName),
	}
	return fmt.Sprintf("Streaming Policy (%s)", strings.Join(components, "\n"))
}
