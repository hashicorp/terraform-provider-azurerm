package scheduledqueryrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScheduledQueryRuleId{})
}

var _ resourceids.ResourceId = &ScheduledQueryRuleId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ScheduledQueryRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScheduledQueryRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScheduledQueryRuleIDInsensitively parses 'input' case-insensitively into a ScheduledQueryRuleId
// note: this method should only be used for API response data and not user input
func ParseScheduledQueryRuleIDInsensitively(input string) (*ScheduledQueryRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScheduledQueryRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScheduledQueryRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScheduledQueryRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ScheduledQueryRuleName, ok = input.Parsed["scheduledQueryRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scheduledQueryRuleName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("scheduledQueryRuleName", "scheduledQueryRuleName"),
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
