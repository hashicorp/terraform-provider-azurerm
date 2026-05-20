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
	recaser.RegisterResourceId(&HubId{})
}

var _ resourceids.ResourceId = &HubId{}

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
	parser := resourceids.NewParserFromResourceIdType(&HubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHubIDInsensitively parses 'input' case-insensitively into a HubId
// note: this method should only be used for API response data and not user input
func ParseHubIDInsensitively(input string) (*HubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HubId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HubName, ok = input.Parsed["hubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hubName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubName"),
		resourceids.StaticSegment("staticHubs", "hubs", "hubs"),
		resourceids.UserSpecifiedSegment("hubName", "hubName"),
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
