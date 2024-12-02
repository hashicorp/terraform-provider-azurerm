package routefilterrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteFilterRuleId{})
}

var _ resourceids.ResourceId = &RouteFilterRuleId{}

// RouteFilterRuleId is a struct representing the Resource ID for a Route Filter Rule
type RouteFilterRuleId struct {
	SubscriptionId      string
	ResourceGroupName   string
	RouteFilterName     string
	RouteFilterRuleName string
}

// NewRouteFilterRuleID returns a new RouteFilterRuleId struct
func NewRouteFilterRuleID(subscriptionId string, resourceGroupName string, routeFilterName string, routeFilterRuleName string) RouteFilterRuleId {
	return RouteFilterRuleId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		RouteFilterName:     routeFilterName,
		RouteFilterRuleName: routeFilterRuleName,
	}
}

// ParseRouteFilterRuleID parses 'input' into a RouteFilterRuleId
func ParseRouteFilterRuleID(input string) (*RouteFilterRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteFilterRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteFilterRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteFilterRuleIDInsensitively parses 'input' case-insensitively into a RouteFilterRuleId
// note: this method should only be used for API response data and not user input
func ParseRouteFilterRuleIDInsensitively(input string) (*RouteFilterRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteFilterRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteFilterRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteFilterRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RouteFilterName, ok = input.Parsed["routeFilterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeFilterName", input)
	}

	if id.RouteFilterRuleName, ok = input.Parsed["routeFilterRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeFilterRuleName", input)
	}

	return nil
}

// ValidateRouteFilterRuleID checks that 'input' can be parsed as a Route Filter Rule ID
func ValidateRouteFilterRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteFilterRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Filter Rule ID
func (id RouteFilterRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeFilters/%s/routeFilterRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RouteFilterName, id.RouteFilterRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Filter Rule ID
func (id RouteFilterRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticRouteFilters", "routeFilters", "routeFilters"),
		resourceids.UserSpecifiedSegment("routeFilterName", "routeFilterName"),
		resourceids.StaticSegment("staticRouteFilterRules", "routeFilterRules", "routeFilterRules"),
		resourceids.UserSpecifiedSegment("routeFilterRuleName", "routeFilterRuleName"),
	}
}

// String returns a human-readable description of this Route Filter Rule ID
func (id RouteFilterRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Route Filter Name: %q", id.RouteFilterName),
		fmt.Sprintf("Route Filter Rule Name: %q", id.RouteFilterRuleName),
	}
	return fmt.Sprintf("Route Filter Rule (%s)", strings.Join(components, "\n"))
}
