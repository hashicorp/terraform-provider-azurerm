package alertprocessingrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ActionRuleId{})
}

var _ resourceids.ResourceId = &ActionRuleId{}

// ActionRuleId is a struct representing the Resource ID for a Action Rule
type ActionRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	ActionRuleName    string
}

// NewActionRuleID returns a new ActionRuleId struct
func NewActionRuleID(subscriptionId string, resourceGroupName string, actionRuleName string) ActionRuleId {
	return ActionRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ActionRuleName:    actionRuleName,
	}
}

// ParseActionRuleID parses 'input' into a ActionRuleId
func ParseActionRuleID(input string) (*ActionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActionRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActionRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseActionRuleIDInsensitively parses 'input' case-insensitively into a ActionRuleId
// note: this method should only be used for API response data and not user input
func ParseActionRuleIDInsensitively(input string) (*ActionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActionRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActionRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ActionRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ActionRuleName, ok = input.Parsed["actionRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "actionRuleName", input)
	}

	return nil
}

// ValidateActionRuleID checks that 'input' can be parsed as a Action Rule ID
func ValidateActionRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActionRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Action Rule ID
func (id ActionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AlertsManagement/actionRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ActionRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Action Rule ID
func (id ActionRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAlertsManagement", "Microsoft.AlertsManagement", "Microsoft.AlertsManagement"),
		resourceids.StaticSegment("staticActionRules", "actionRules", "actionRules"),
		resourceids.UserSpecifiedSegment("actionRuleName", "actionRuleName"),
	}
}

// String returns a human-readable description of this Action Rule ID
func (id ActionRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Action Rule Name: %q", id.ActionRuleName),
	}
	return fmt.Sprintf("Action Rule (%s)", strings.Join(components, "\n"))
}
