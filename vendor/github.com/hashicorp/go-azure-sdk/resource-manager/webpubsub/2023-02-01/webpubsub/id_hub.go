package webpubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = HubId{}

// HubId is a struct representing the Resource ID for a Hub
type HubId struct {
	SubscriptionId    string
	ResourceGroupName string
	WebPubSubName     string
	HubName           string
}

// NewHubID returns a new HubId struct
func NewHubID(subscriptionId string, resourceGroupName string, webPubSubName string, hubName string) HubId {
	return HubId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WebPubSubName:     webPubSubName,
		HubName:           hubName,
	}
}

// ParseHubID parses 'input' into a HubId
func ParseHubID(input string) (*HubId, error) {
	parser := resourceids.NewParserFromResourceIdType(HubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	if id.HubName, ok = parsed.Parsed["hubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hubName", *parsed)
	}

	return &id, nil
}

// ParseHubIDInsensitively parses 'input' case-insensitively into a HubId
// note: this method should only be used for API response data and not user input
func ParseHubIDInsensitively(input string) (*HubId, error) {
	parser := resourceids.NewParserFromResourceIdType(HubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	if id.HubName, ok = parsed.Parsed["hubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hubName", *parsed)
	}

	return &id, nil
}

// ValidateHubID checks that 'input' can be parsed as a Hub ID
func ValidateHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hub ID
func (id HubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/hubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName, id.HubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hub ID
func (id HubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticWebPubSub", "webPubSub", "webPubSub"),
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubValue"),
		resourceids.StaticSegment("staticHubs", "hubs", "hubs"),
		resourceids.UserSpecifiedSegment("hubName", "hubValue"),
	}
}

// String returns a human-readable description of this Hub ID
func (id HubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Web Pub Sub Name: %q", id.WebPubSubName),
		fmt.Sprintf("Hub Name: %q", id.HubName),
	}
	return fmt.Sprintf("Hub (%s)", strings.Join(components, "\n"))
}
