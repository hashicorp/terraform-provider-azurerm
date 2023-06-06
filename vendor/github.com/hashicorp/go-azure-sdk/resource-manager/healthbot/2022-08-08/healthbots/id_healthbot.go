package healthbots

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = HealthBotId{}

// HealthBotId is a struct representing the Resource ID for a Health Bot
type HealthBotId struct {
	SubscriptionId    string
	ResourceGroupName string
	HealthBotName     string
}

// NewHealthBotID returns a new HealthBotId struct
func NewHealthBotID(subscriptionId string, resourceGroupName string, healthBotName string) HealthBotId {
	return HealthBotId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HealthBotName:     healthBotName,
	}
}

// ParseHealthBotID parses 'input' into a HealthBotId
func ParseHealthBotID(input string) (*HealthBotId, error) {
	parser := resourceids.NewParserFromResourceIdType(HealthBotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HealthBotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HealthBotName, ok = parsed.Parsed["healthBotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "healthBotName", *parsed)
	}

	return &id, nil
}

// ParseHealthBotIDInsensitively parses 'input' case-insensitively into a HealthBotId
// note: this method should only be used for API response data and not user input
func ParseHealthBotIDInsensitively(input string) (*HealthBotId, error) {
	parser := resourceids.NewParserFromResourceIdType(HealthBotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HealthBotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HealthBotName, ok = parsed.Parsed["healthBotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "healthBotName", *parsed)
	}

	return &id, nil
}

// ValidateHealthBotID checks that 'input' can be parsed as a Health Bot ID
func ValidateHealthBotID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHealthBotID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Health Bot ID
func (id HealthBotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthBot/healthBots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HealthBotName)
}

// Segments returns a slice of Resource ID Segments which comprise this Health Bot ID
func (id HealthBotId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHealthBot", "Microsoft.HealthBot", "Microsoft.HealthBot"),
		resourceids.StaticSegment("staticHealthBots", "healthBots", "healthBots"),
		resourceids.UserSpecifiedSegment("healthBotName", "healthBotValue"),
	}
}

// String returns a human-readable description of this Health Bot ID
func (id HealthBotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Health Bot Name: %q", id.HealthBotName),
	}
	return fmt.Sprintf("Health Bot (%s)", strings.Join(components, "\n"))
}
