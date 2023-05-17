package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RulesEngineId{}

// RulesEngineId is a struct representing the Resource ID for a Rules Engine
type RulesEngineId struct {
	SubscriptionId    string
	ResourceGroupName string
	FrontDoorName     string
	RulesEngineName   string
}

// NewRulesEngineID returns a new RulesEngineId struct
func NewRulesEngineID(subscriptionId string, resourceGroupName string, frontDoorName string, rulesEngineName string) RulesEngineId {
	return RulesEngineId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FrontDoorName:     frontDoorName,
		RulesEngineName:   rulesEngineName,
	}
}

// ParseRulesEngineID parses 'input' into a RulesEngineId
func ParseRulesEngineID(input string) (*RulesEngineId, error) {
	parser := resourceids.NewParserFromResourceIdType(RulesEngineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RulesEngineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	if id.RulesEngineName, ok = parsed.Parsed["rulesEngineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "rulesEngineName", *parsed)
	}

	return &id, nil
}

// ParseRulesEngineIDInsensitively parses 'input' case-insensitively into a RulesEngineId
// note: this method should only be used for API response data and not user input
func ParseRulesEngineIDInsensitively(input string) (*RulesEngineId, error) {
	parser := resourceids.NewParserFromResourceIdType(RulesEngineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RulesEngineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	if id.RulesEngineName, ok = parsed.Parsed["rulesEngineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "rulesEngineName", *parsed)
	}

	return &id, nil
}

// ValidateRulesEngineID checks that 'input' can be parsed as a Rules Engine ID
func ValidateRulesEngineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRulesEngineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rules Engine ID
func (id RulesEngineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/rulesEngines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FrontDoorName, id.RulesEngineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rules Engine ID
func (id RulesEngineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFrontDoors", "frontDoors", "frontDoors"),
		resourceids.UserSpecifiedSegment("frontDoorName", "frontDoorValue"),
		resourceids.StaticSegment("staticRulesEngines", "rulesEngines", "rulesEngines"),
		resourceids.UserSpecifiedSegment("rulesEngineName", "rulesEngineValue"),
	}
}

// String returns a human-readable description of this Rules Engine ID
func (id RulesEngineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Front Door Name: %q", id.FrontDoorName),
		fmt.Sprintf("Rules Engine Name: %q", id.RulesEngineName),
	}
	return fmt.Sprintf("Rules Engine (%s)", strings.Join(components, "\n"))
}
