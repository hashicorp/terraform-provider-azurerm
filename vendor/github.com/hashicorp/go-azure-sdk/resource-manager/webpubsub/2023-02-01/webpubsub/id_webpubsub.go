package webpubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WebPubSubId{}

// WebPubSubId is a struct representing the Resource ID for a Web Pub Sub
type WebPubSubId struct {
	SubscriptionId    string
	ResourceGroupName string
	WebPubSubName     string
}

// NewWebPubSubID returns a new WebPubSubId struct
func NewWebPubSubID(subscriptionId string, resourceGroupName string, webPubSubName string) WebPubSubId {
	return WebPubSubId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WebPubSubName:     webPubSubName,
	}
}

// ParseWebPubSubID parses 'input' into a WebPubSubId
func ParseWebPubSubID(input string) (*WebPubSubId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebPubSubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebPubSubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	return &id, nil
}

// ParseWebPubSubIDInsensitively parses 'input' case-insensitively into a WebPubSubId
// note: this method should only be used for API response data and not user input
func ParseWebPubSubIDInsensitively(input string) (*WebPubSubId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebPubSubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebPubSubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebPubSubName, ok = parsed.Parsed["webPubSubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", *parsed)
	}

	return &id, nil
}

// ValidateWebPubSubID checks that 'input' can be parsed as a Web Pub Sub ID
func ValidateWebPubSubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebPubSubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Web Pub Sub ID
func (id WebPubSubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Web Pub Sub ID
func (id WebPubSubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticWebPubSub", "webPubSub", "webPubSub"),
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubValue"),
	}
}

// String returns a human-readable description of this Web Pub Sub ID
func (id WebPubSubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Web Pub Sub Name: %q", id.WebPubSubName),
	}
	return fmt.Sprintf("Web Pub Sub (%s)", strings.Join(components, "\n"))
}
