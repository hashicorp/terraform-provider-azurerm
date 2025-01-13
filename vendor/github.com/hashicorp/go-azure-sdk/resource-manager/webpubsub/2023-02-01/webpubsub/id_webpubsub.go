package webpubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WebPubSubId{})
}

var _ resourceids.ResourceId = &WebPubSubId{}

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
	parser := resourceids.NewParserFromResourceIdType(&WebPubSubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebPubSubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWebPubSubIDInsensitively parses 'input' case-insensitively into a WebPubSubId
// note: this method should only be used for API response data and not user input
func ParseWebPubSubIDInsensitively(input string) (*WebPubSubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WebPubSubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebPubSubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WebPubSubId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WebPubSubName, ok = input.Parsed["webPubSubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubName"),
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
