// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &BotServiceId{}

// BotServiceId is a struct representing the Resource ID for a Bot Service
type BotServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	BotServiceName    string
}

// NewBotServiceID returns a new BotServiceId struct
func NewBotServiceID(subscriptionId string, resourceGroupName string, botServiceName string) BotServiceId {
	return BotServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BotServiceName:    botServiceName,
	}
}

// ParseBotServiceID parses 'input' into a BotServiceId
func ParseBotServiceID(input string) (*BotServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BotServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BotServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBotServiceIDInsensitively parses 'input' case-insensitively into a BotServiceId
// note: this method should only be used for API response data and not user input
func ParseBotServiceIDInsensitively(input string) (*BotServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BotServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BotServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BotServiceId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.BotServiceName, ok = input.Parsed["botServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "botServiceName", input)
	}

	return nil
}

// ValidateBotServiceID checks that 'input' can be parsed as a Bot Service ID
func ValidateBotServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBotServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bot Service ID
func (id BotServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.BotService/botServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BotServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Bot Service ID
func (id BotServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBotService", "Microsoft.BotService", "Microsoft.BotService"),
		resourceids.StaticSegment("staticBotServices", "botServices", "botServices"),
		resourceids.UserSpecifiedSegment("botServiceName", "botServiceValue"),
	}
}

// String returns a human-readable description of this Bot Service ID
func (id BotServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bot Service Name: %q", id.BotServiceName),
	}
	return fmt.Sprintf("Bot Service (%s)", strings.Join(components, "\n"))
}
