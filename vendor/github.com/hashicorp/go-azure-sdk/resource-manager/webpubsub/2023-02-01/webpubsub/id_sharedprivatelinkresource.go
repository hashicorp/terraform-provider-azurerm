package webpubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SharedPrivateLinkResourceId{}

// SharedPrivateLinkResourceId is a struct representing the Resource ID for a Shared Private Link Resource
type SharedPrivateLinkResourceId struct {
	SubscriptionId                string
	ResourceGroupName             string
	WebPubSubName                 string
	SharedPrivateLinkResourceName string
}

// NewSharedPrivateLinkResourceID returns a new SharedPrivateLinkResourceId struct
func NewSharedPrivateLinkResourceID(subscriptionId string, resourceGroupName string, webPubSubName string, sharedPrivateLinkResourceName string) SharedPrivateLinkResourceId {
	return SharedPrivateLinkResourceId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		WebPubSubName:                 webPubSubName,
		SharedPrivateLinkResourceName: sharedPrivateLinkResourceName,
	}
}

// ParseSharedPrivateLinkResourceID parses 'input' into a SharedPrivateLinkResourceId
func ParseSharedPrivateLinkResourceID(input string) (*SharedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SharedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SharedPrivateLinkResourceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	if id.SharedPrivateLinkResourceName, ok = parsed.Parsed["sharedPrivateLinkResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sharedPrivateLinkResourceName", *parsed)
	}

	return &id, nil
}

// ParseSharedPrivateLinkResourceIDInsensitively parses 'input' case-insensitively into a SharedPrivateLinkResourceId
// note: this method should only be used for API response data and not user input
func ParseSharedPrivateLinkResourceIDInsensitively(input string) (*SharedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SharedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SharedPrivateLinkResourceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	if id.SharedPrivateLinkResourceName, ok = parsed.Parsed["sharedPrivateLinkResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sharedPrivateLinkResourceName", *parsed)
	}

	return &id, nil
}

// ValidateSharedPrivateLinkResourceID checks that 'input' can be parsed as a Shared Private Link Resource ID
func ValidateSharedPrivateLinkResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSharedPrivateLinkResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Shared Private Link Resource ID
func (id SharedPrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/sharedPrivateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName, id.SharedPrivateLinkResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Shared Private Link Resource ID
func (id SharedPrivateLinkResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticWebPubSub", "webPubSub", "webPubSub"),
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubValue"),
		resourceids.StaticSegment("staticSharedPrivateLinkResources", "sharedPrivateLinkResources", "sharedPrivateLinkResources"),
		resourceids.UserSpecifiedSegment("sharedPrivateLinkResourceName", "sharedPrivateLinkResourceValue"),
	}
}

// String returns a human-readable description of this Shared Private Link Resource ID
func (id SharedPrivateLinkResourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Web Pub Sub Name: %q", id.WebPubSubName),
		fmt.Sprintf("Shared Private Link Resource Name: %q", id.SharedPrivateLinkResourceName),
	}
	return fmt.Sprintf("Shared Private Link Resource (%s)", strings.Join(components, "\n"))
}
