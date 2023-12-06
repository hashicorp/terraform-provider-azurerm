package liveevents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LiveEventOperationId{}

// LiveEventOperationId is a struct representing the Resource ID for a Live Event Operation
type LiveEventOperationId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	OperationId       string
}

// NewLiveEventOperationID returns a new LiveEventOperationId struct
func NewLiveEventOperationID(subscriptionId string, resourceGroupName string, mediaServiceName string, operationId string) LiveEventOperationId {
	return LiveEventOperationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		OperationId:       operationId,
	}
}

// ParseLiveEventOperationID parses 'input' into a LiveEventOperationId
func ParseLiveEventOperationID(input string) (*LiveEventOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveEventOperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LiveEventOperationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLiveEventOperationIDInsensitively parses 'input' case-insensitively into a LiveEventOperationId
// note: this method should only be used for API response data and not user input
func ParseLiveEventOperationIDInsensitively(input string) (*LiveEventOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LiveEventOperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LiveEventOperationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LiveEventOperationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MediaServiceName, ok = input.Parsed["mediaServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", input)
	}

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	return nil
}

// ValidateLiveEventOperationID checks that 'input' can be parsed as a Live Event Operation ID
func ValidateLiveEventOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLiveEventOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Live Event Operation ID
func (id LiveEventOperationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/liveEventOperations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Live Event Operation ID
func (id LiveEventOperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticLiveEventOperations", "liveEventOperations", "liveEventOperations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Live Event Operation ID
func (id LiveEventOperationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Live Event Operation (%s)", strings.Join(components, "\n"))
}
