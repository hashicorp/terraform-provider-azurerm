package scheduledqueryrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScheduledQueryRuleId{}

// ScheduledQueryRuleId is a struct representing the Resource ID for a Scheduled Query Rule
type ScheduledQueryRuleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ScheduledQueryRuleName string
}

// NewScheduledQueryRuleID returns a new ScheduledQueryRuleId struct
func NewScheduledQueryRuleID(subscriptionId string, resourceGroupName string, scheduledQueryRuleName string) ScheduledQueryRuleId {
	return ScheduledQueryRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ScheduledQueryRuleName: scheduledQueryRuleName,
	}
}

// ParseScheduledQueryRuleID parses 'input' into a ScheduledQueryRuleId
func ParseScheduledQueryRuleID(input string) (*ScheduledQueryRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduledQueryRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduledQueryRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ScheduledQueryRuleName, ok = parsed.Parsed["scheduledQueryRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scheduledQueryRuleName", *parsed)
	}

	return &id, nil
}

// ParseScheduledQueryRuleIDInsensitively parses 'input' case-insensitively into a ScheduledQueryRuleId
// note: this method should only be used for API response data and not user input
func ParseScheduledQueryRuleIDInsensitively(input string) (*ScheduledQueryRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduledQueryRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduledQueryRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ScheduledQueryRuleName, ok = parsed.Parsed["scheduledQueryRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scheduledQueryRuleName", *parsed)
	}

	return &id, nil
}

// ValidateScheduledQueryRuleID checks that 'input' can be parsed as a Scheduled Query Rule ID
func ValidateScheduledQueryRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScheduledQueryRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scheduled Query Rule ID
func (id ScheduledQueryRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/scheduledQueryRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ScheduledQueryRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scheduled Query Rule ID
func (id ScheduledQueryRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticScheduledQueryRules", "scheduledQueryRules", "scheduledQueryRules"),
		resourceids.UserSpecifiedSegment("scheduledQueryRuleName", "scheduledQueryRuleValue"),
	}
}

// String returns a human-readable description of this Scheduled Query Rule ID
func (id ScheduledQueryRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Scheduled Query Rule Name: %q", id.ScheduledQueryRuleName),
	}
	return fmt.Sprintf("Scheduled Query Rule (%s)", strings.Join(components, "\n"))
}
